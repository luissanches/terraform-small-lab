package kt_session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRead,
		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "->> Starting Instances Read...")
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 25 * time.Second}
	dataProvider := m.(DataProvider)

	url := fmt.Sprintf("http://localhost:5000/instances?username=%v", dataProvider.user_name)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-type", "application/json")
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	instance := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("instances", instance); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	tflog.Debug(ctx, "->> Read Instances Completed Successfully...")
	return diags
}
