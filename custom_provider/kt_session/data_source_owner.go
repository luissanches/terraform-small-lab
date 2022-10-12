package kt_session

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOwner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOwnerRead,
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "->> Starting owner Read...")
	var diags diag.Diagnostics

	client := &http.Client{Timeout: 25 * time.Second}

	url := "http://localhost:5000/owner"
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

	
	var instance map[string]interface{} = nil

	err = json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return diag.FromErr(err)
	}

	instances := make([]map[string]interface{}, 0)
	instances = append(instances, instance)
	if err := d.Set("owner", instances); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	tflog.Debug(ctx, "->> Read owner Completed Successfully...")
	return diags
}
