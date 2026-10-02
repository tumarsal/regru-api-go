package main

import (
	"context"
	"fmt"

	regru "github.com/tumarsal/regru-api-go"
	"github.com/spf13/cobra"
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

func runZoneSync(cmd *cobra.Command, args []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	prune, _ := cmd.Flags().GetBool("prune")
	domain := args[0]

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	result, err := regru.SyncZone(context.Background(), client, domain, filePath, regru.SyncOptions{
		DryRun: dryRun,
		Prune:  prune,
	})
	if result != nil {
		if result.ExportYAML != "" {
			fmt.Print(result.ExportYAML)
		}
		for _, c := range result.Changes {
			fmt.Println(c.Line)
		}
		if result.Message != "" {
			fmt.Println(result.Message)
		}
	}
	return err
}
