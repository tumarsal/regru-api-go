package main

import (
	"context"
	"fmt"
	"strings"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/cobra"
)

var zoneAddCmd = &cobra.Command{
	Use:   "add [domain]",
	Short: "Добавить DNS-запись в зону домена",
	Long: `Добавляет запись через REG.API zone/add_*.

Типы: A, AAAA, CNAME, MX, NS, TXT, SRV, CAA, HTTPS.

Примеры:
  regru zone add example.ru --type A --subdomain www --content 1.2.3.4
  regru zone add example.ru --type AAAA --subdomain @ --content 2001:db8::1
  regru zone add example.ru --type CNAME --subdomain mail --canonical-name mx.example.ru
  regru zone add example.ru --type MX --subdomain @ --content mail.example.ru --priority 10
  regru zone add example.ru --type TXT --subdomain @ --text "v=spf1 include:_spf.example.ru ~all"
  regru zone add example.ru --type SRV --service _sip._udp --target sip.example.ru --port 5060`,
	Args: cobra.ExactArgs(1),
	RunE: runZoneAdd,
}

func init() {
	zoneAddCmd.Flags().String("type", "", "тип записи (обязательно)")
	zoneAddCmd.Flags().String("subdomain", "@", "поддомен (@ — корень, * — wildcard)")
	zoneAddCmd.Flags().String("content", "", "содержимое (IP, текст, mail_server, target для простых типов)")
	zoneAddCmd.Flags().String("priority", "", "приоритет (MX, SRV, HTTPS)")
	zoneAddCmd.Flags().String("canonical-name", "", "целевой домен для CNAME")
	zoneAddCmd.Flags().String("dns-server", "", "DNS-сервер для NS")
	zoneAddCmd.Flags().String("record-number", "", "номер NS-записи")
	zoneAddCmd.Flags().String("service", "", "сервис для SRV (например _sip._udp)")
	zoneAddCmd.Flags().String("target", "", "target для SRV/HTTPS")
	zoneAddCmd.Flags().String("port", "", "порт для SRV")
	zoneAddCmd.Flags().String("weight", "", "weight для SRV")
	zoneAddCmd.Flags().String("text", "", "текст для TXT")
	zoneAddCmd.Flags().Int("flags", 0, "флаги для CAA (0 или 128)")
	zoneAddCmd.Flags().String("tag", "", "тег CAA: issue, issuewild, iodef")
	zoneAddCmd.Flags().String("value", "", "значение для CAA/HTTPS")
	_ = zoneAddCmd.MarkFlagRequired("type")
	zoneCmd.AddCommand(zoneAddCmd)
}

func runZoneAdd(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	rec, err := dnsRecordFromFlags(cmd)
	if err != nil {
		return err
	}

	if err := client.Zone().AddRecord(context.Background(), args[0], rec); err != nil {
		return err
	}

	fmt.Printf("запись %s для %s добавлена\n", strings.ToUpper(rec.RecordType), args[0])
	return nil
}

func dnsRecordFromFlags(cmd *cobra.Command) (regru.DNSRecord, error) {
	typ, _ := cmd.Flags().GetString("type")
	sub, _ := cmd.Flags().GetString("subdomain")
	content, _ := cmd.Flags().GetString("content")
	priority, _ := cmd.Flags().GetString("priority")
	cname, _ := cmd.Flags().GetString("canonical-name")
	dnsServer, _ := cmd.Flags().GetString("dns-server")
	recordNum, _ := cmd.Flags().GetString("record-number")
	service, _ := cmd.Flags().GetString("service")
	target, _ := cmd.Flags().GetString("target")
	port, _ := cmd.Flags().GetString("port")
	weight, _ := cmd.Flags().GetString("weight")
	text, _ := cmd.Flags().GetString("text")
	flags, _ := cmd.Flags().GetInt("flags")
	tag, _ := cmd.Flags().GetString("tag")
	value, _ := cmd.Flags().GetString("value")

	rec := regru.DNSRecord{
		RecordType:    typ,
		SubDomain:     sub,
		Content:       content,
		Priority:      priority,
		CanonicalName: cname,
		DNSServer:     dnsServer,
		RecordNumber:  recordNum,
		Service:       service,
		Target:        target,
		Port:          port,
		Weight:        weight,
		Text:          text,
		Flags:         flags,
		Tag:           tag,
		Value:         value,
	}

	if strings.ToUpper(typ) == "HTTPS" && priority == "" {
		rec.Priority = "0"
	}

	return rec, nil
}
