package commands

import (
	"errors"
	"io"

	"github.com/bryanl/doit"
	"github.com/bryanl/doit/do"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// FloatingIP creates the command heirarchy for floating ips.
func FloatingIP() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "floating-ip",
		Short:   "floating IP commands",
		Long:    "floating-ip is used to access commands on floating IPs",
		Aliases: []string{"fip"},
	}

	cmdFloatingIPCreate := cmdBuilder(cmd, RunFloatingIPCreate, "create", "create a floating IP", writer,
		aliasOpt("c"), displayerType(&floatingIP{}))
	addStringFlag(cmdFloatingIPCreate, doit.ArgRegionSlug, "", "Region where to create the floating IP.", requiredOpt())
	addIntFlag(cmdFloatingIPCreate, doit.ArgDropletID, 0, "ID of the droplet to assign the IP to. (Optional)")

	cmdBuilder(cmd, RunFloatingIPGet, "get <floating-ip>", "get the details of a floating IP", writer,
		aliasOpt("g"), displayerType(&floatingIP{}))

	cmdBuilder(cmd, RunFloatingIPDelete, "delete <floating-ip>", "delete a floating IP address", writer, aliasOpt("d"))

	cmdFloatingIPList := cmdBuilder(cmd, RunFloatingIPList, "list", "list all floating IP addresses", writer,
		aliasOpt("ls"), displayerType(&floatingIP{}))
	addStringFlag(cmdFloatingIPList, doit.ArgRegionSlug, "", "Floating IP region")

	return cmd
}

// RunFloatingIPCreate runs floating IP create.
func RunFloatingIPCreate(ns string, config doit.Config, out io.Writer, args []string) error {
	client := config.GetGodoClient()
	fis := do.NewFloatingIPsService(client)

	region, err := config.GetString(ns, doit.ArgRegionSlug)
	if err != nil {
		return err
	}

	dropletID, err := config.GetInt(ns, doit.ArgDropletID)
	if err != nil {
		return err
	}

	req := &godo.FloatingIPCreateRequest{
		Region:    region,
		DropletID: dropletID,
	}

	ip, err := fis.Create(req)
	if err != nil {
		return err
	}

	dc := &displayer{
		ns:     ns,
		config: config,
		item:   &floatingIP{floatingIPs: do.FloatingIPs{*ip}},
		out:    out,
	}
	return dc.Display()
}

// RunFloatingIPGet retrieves a floating IP's details.
func RunFloatingIPGet(ns string, config doit.Config, out io.Writer, args []string) error {
	client := config.GetGodoClient()
	fis := do.NewFloatingIPsService(client)

	if len(args) != 1 {
		return doit.NewMissingArgsErr(ns)
	}

	ip := args[0]

	if len(ip) < 1 {
		return errors.New("invalid ip address")
	}

	fip, err := fis.Get(ip)
	if err != nil {
		return err
	}

	dc := &displayer{
		ns:     ns,
		config: config,
		item:   &floatingIP{floatingIPs: do.FloatingIPs{*fip}},
		out:    out,
	}

	return dc.Display()
}

// RunFloatingIPDelete runs floating IP delete.
func RunFloatingIPDelete(ns string, config doit.Config, out io.Writer, args []string) error {
	client := config.GetGodoClient()
	ds := do.NewFloatingIPsService(client)

	if len(args) != 1 {
		return doit.NewMissingArgsErr(ns)
	}

	ip := args[0]

	return ds.Delete(ip)
}

// RunFloatingIPList runs floating IP create.
func RunFloatingIPList(ns string, config doit.Config, out io.Writer, args []string) error {
	client := config.GetGodoClient()
	fis := do.NewFloatingIPsService(client)

	region, err := config.GetString(ns, doit.ArgRegionSlug)
	if err != nil {
		return err
	}

	list, err := fis.List()
	if err != nil {
		return err
	}

	fips := &floatingIP{floatingIPs: do.FloatingIPs{}}
	for _, fip := range list {
		var skip bool
		if region != "" && region != fip.Region.Slug {
			skip = true
		}

		if !skip {
			fips.floatingIPs = append(fips.floatingIPs, fip)
		}
	}

	dc := &displayer{
		ns:     ns,
		config: config,
		item:   fips,
		out:    out,
	}

	return dc.Display()
}