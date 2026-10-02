package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var zoneSyncCmd = &cobra.Command{
	Use:   "sync [domain]",
	Short: "Синхронизировать DNS-записи зоны с YAML-файлом",
	Long: `Приводит зону домена к состоянию из файла: добавляет недостающие записи
и (по умолчанию) удаляет лишние.

Если файл по пути -f не существует, записи загружаются из API и сохраняются в YAML
(синхронизация при этом не выполняется).

Формат файла (поля записей совпадают с выводом «zone list»):

  records:
    - type: A
      subdomain: www
      content: 1.2.3.4
    - type: MX
      subdomain: "@"
      priority: "10"
      content: mail.example.ru

Для типов CNAME, NS, TXT, SRV, CAA, HTTPS можно задавать те же поля, что и у «zone add»
(canonical_name, dns_server, service, target, text, tag, value и т.д.).

Примеры:
  regru zone sync example.ru -f dns.yaml
  regru zone sync example.ru -f dns.yaml --dry-run
  regru zone sync example.ru -f dns.yaml --prune=false`,
	Args: cobra.ExactArgs(1),
	RunE: runZoneSync,
}

func init() {
	zoneSyncCmd.Flags().StringP("file", "f", "", "путь к YAML-файлу (обязательно)")
	zoneSyncCmd.Flags().Bool("dry-run", false, "только показать план изменений")
	zoneSyncCmd.Flags().Bool("prune", true, "удалять записи, которых нет в файле")
	_ = zoneSyncCmd.MarkFlagRequired("file")
	zoneCmd.AddCommand(zoneSyncCmd)
}

type zoneYAMLFile struct {
	Domain  string           `yaml:"domain,omitempty"`
	Records []zoneYAMLRecord `yaml:"records"`
}

type zoneYAMLRecord struct {
	Type          string `yaml:"type"`
	Subdomain     string `yaml:"subdomain,omitempty"`
	Content       string `yaml:"content,omitempty"`
	Priority      string `yaml:"priority,omitempty"`
	CanonicalName string `yaml:"canonical_name,omitempty"`
	DNSServer     string `yaml:"dns_server,omitempty"`
	RecordNumber  string `yaml:"record_number,omitempty"`
	Service       string `yaml:"service,omitempty"`
	Target        string `yaml:"target,omitempty"`
	Port          string `yaml:"port,omitempty"`
	Weight        string `yaml:"weight,omitempty"`
	Text          string `yaml:"text,omitempty"`
	Flags         int    `yaml:"flags,omitempty"`
	Tag           string `yaml:"tag,omitempty"`
	Value         string `yaml:"value,omitempty"`
}

func runZoneSync(cmd *cobra.Command, args []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	prune, _ := cmd.Flags().GetBool("prune")
	domain := args[0]

	client, err := newAPIClient()
	if err != nil {
		return err
	}
	ctx := context.Background()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return exportZoneToYAML(ctx, client, domain, filePath, dryRun)
		}
		return fmt.Errorf("чтение файла: %w", err)
	}

	var spec zoneYAMLFile
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("разбор YAML: %w", err)
	}
	if len(spec.Records) == 0 {
		return fmt.Errorf("в файле нет записей (ключ records)")
	}

	if spec.Domain != "" && spec.Domain != domain {
		return fmt.Errorf("домен в аргументе (%s) не совпадает с domain в файле (%s)", domain, spec.Domain)
	}

	desired, err := dnsRecordsFromYAML(spec.Records)
	if err != nil {
		return err
	}

	current, err := client.Zone().GetResourceRecords(ctx, domain)
	if err != nil {
		return err
	}

	plan := regru.PlanZoneSync(desired, current, prune)

	if len(plan.ToAdd) == 0 && len(plan.ToRemove) == 0 {
		fmt.Println("изменений нет")
		return nil
	}

	for _, r := range plan.ToRemove {
		line := fmt.Sprintf("- %s %s %s", r.Rectype, r.Subname, r.Content)
		if p := r.Priority.String(); p != "" {
			line += " (prio " + p + ")"
		}
		fmt.Println(line)
	}
	for _, d := range plan.ToAdd {
		key := regru.DesiredRecordKey(d)
		parts := strings.Split(key, "\x00")
		line := fmt.Sprintf("+ %s %s %s", parts[0], parts[1], parts[2])
		if len(parts) > 3 && parts[3] != "" {
			line += " (prio " + parts[3] + ")"
		}
		fmt.Println(line)
	}

	if dryRun {
		fmt.Println("(dry-run: изменения не применены)")
		return nil
	}

	for _, r := range plan.ToRemove {
		if err := client.Zone().RemoveRecord(ctx, domain, r.Subname, r.Rectype, r.Content, r.Priority.String()); err != nil {
			return fmt.Errorf("удаление %s %s: %w", r.Rectype, r.Subname, err)
		}
	}
	for _, d := range plan.ToAdd {
		if err := client.Zone().AddRecord(ctx, domain, d); err != nil {
			return fmt.Errorf("добавление %s %s: %w", d.RecordType, d.SubDomain, err)
		}
	}

	fmt.Printf("синхронизация %s завершена: удалено %d, добавлено %d\n", domain, len(plan.ToRemove), len(plan.ToAdd))
	return nil
}

func exportZoneToYAML(ctx context.Context, client *regru.Client, domain, filePath string, dryRun bool) error {
	records, err := client.Zone().GetResourceRecords(ctx, domain)
	if err != nil {
		return err
	}

	spec := zoneYAMLFile{
		Domain:  domain,
		Records: yamlRecordsFromAPI(records),
	}
	data, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("формирование YAML: %w", err)
	}

	if dryRun {
		fmt.Print(string(data))
		fmt.Printf("(dry-run: файл %s не записан)\n", filePath)
		return nil
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("запись файла: %w", err)
	}
	fmt.Printf("файл %s создан: %d записей из API\n", filePath, len(spec.Records))
	return nil
}

func yamlRecordsFromAPI(records []regru.ResourceRecord) []zoneYAMLRecord {
	out := make([]zoneYAMLRecord, 0, len(records))
	for _, r := range records {
		row := zoneYAMLRecord{
			Type:      r.Rectype,
			Subdomain: r.Subname,
			Content:   r.Content,
		}
		if p := r.Priority.String(); p != "" {
			row.Priority = p
		}
		out = append(out, row)
	}
	return out
}

func dnsRecordsFromYAML(rows []zoneYAMLRecord) ([]regru.DNSRecord, error) {
	out := make([]regru.DNSRecord, 0, len(rows))
	for i, row := range rows {
		if strings.TrimSpace(row.Type) == "" {
			return nil, fmt.Errorf("запись %d: не указан type", i+1)
		}
		rec := regru.DNSRecord{
			RecordType:    row.Type,
			SubDomain:     row.Subdomain,
			Content:       row.Content,
			Priority:      row.Priority,
			CanonicalName: row.CanonicalName,
			DNSServer:     row.DNSServer,
			RecordNumber:  row.RecordNumber,
			Service:       row.Service,
			Target:        row.Target,
			Port:          row.Port,
			Weight:        row.Weight,
			Text:          row.Text,
			Flags:         row.Flags,
			Tag:           row.Tag,
			Value:         row.Value,
		}
		if strings.ToUpper(strings.TrimSpace(rec.RecordType)) == "HTTPS" && rec.Priority == "" {
			rec.Priority = "0"
		}
		out = append(out, rec)
	}
	return out, nil
}
