package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	regru "github.com/example/regru-api-go"
	"github.com/spf13/cobra"
)

var zoneListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "Показать DNS-записи домена",
	Args:  cobra.ExactArgs(1),
	RunE:  runZoneList,
}

func init() {
	zoneListCmd.Flags().String("type", "", "фильтр по типу (A, AAAA, CNAME, MX, NS, TXT, SRV, CAA, HTTPS)")
	zoneListCmd.Flags().String("subdomain", "", "фильтр по поддомену (@, *, www, ...)")
	zoneListCmd.Flags().String("content", "", "фильтр по содержимому (точное совпадение)")
	zoneCmd.AddCommand(zoneListCmd)
}

func runZoneList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	records, err := client.Zone().GetResourceRecords(context.Background(), args[0])
	if err != nil {
		return err
	}

	filter := regru.RecordFilter{}
	filter.Type, _ = cmd.Flags().GetString("type")
	filter.Subdomain, _ = cmd.Flags().GetString("subdomain")
	filter.Content, _ = cmd.Flags().GetString("content")

	records = regru.FilterRecords(records, filter)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tSUBDOMAIN\tPRIORITY\tCONTENT")
	for _, r := range records {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Rectype, r.Subname, r.Priority, r.Content)
	}
	return w.Flush()
}
