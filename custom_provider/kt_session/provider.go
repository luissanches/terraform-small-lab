package kt_session

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataProvider struct {
	user_name          string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},	
		ResourcesMap: map[string]*schema.Resource{
			"kt-session_instance": resourceInstance(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kt-session_owner":	dataSourceOwner(),
			"kt-session_instances":	dataSourceInstance(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Debug(ctx, "->> Starting Provider Configuration...")

	var diags diag.Diagnostics

	dataProvider := DataProvider{
		user_name:          d.Get("user_name").(string),
	}

	tflog.Debug(ctx, "->> Provider Configuration Completed Successfully...")
	return dataProvider, diags
}
