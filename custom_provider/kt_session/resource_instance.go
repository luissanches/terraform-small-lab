package kt_session

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


type instance struct {
	Type 		string     `json:"type"`
}

type responseData struct {
	Type 					string     `json:"type"`
	Status 				string        `json:"status"`
	ID 						string     `json:"id"`
	ResourceName 	string     `json:"resource_name"`
}

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "->> Instance Create starting...")

	var diags diag.Diagnostics

	instance := instance{
		Type:      d.Get("type").(string),
	}

	requestBody, err := json.Marshal(instance)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error on read API response.",
			Detail:   "Error on read API response Detail.",
		})

		return diags
	}

	url := "http://localhost:5000/instance"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Add("Content-type", "application/json")
	request.Close = true

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error on make API request.",
			Detail:   "Error on make API request Detail.",
		})

		return diags
	}

	client := http.Client{
		Timeout: time.Duration(time.Minute * 5),
	}

	response, err := client.Do(request)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error on read API response.",
			Detail:   "Error on read API response Detail.",
		})

		return diags
	}
	defer response.Body.Close()


	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error on read API response.",
			Detail:   "Error on read API response Detail.",
		})
		return diags
	}

	resultString := string(body)
	tflog.Debug(ctx, fmt.Sprintf("\n\n ####### Response Status Code: %v \n\n\n", response.StatusCode))
	tflog.Debug(ctx, fmt.Sprintf("\n\n ####### Response Payload: %v \n\n\n", resultString))

	var resultData responseData
	if responseErr := json.Unmarshal(body, &resultData); responseErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error to parse response API.",
			Detail:   "Error to parse response API.",
		})
		return diags
	}


	if response.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  resultData.ID,
			Detail:   fmt.Sprintf("Please verify the resource kt-session_instance.\nTraceId: %v.", resultData.ID),
		})
		return diags
	}

	d.SetId(resultData.ID)
	tflog.Debug(ctx, "->> Instance Create Completed Successfully...")
	return diags
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "->> Action Read Starting...")
	var diags diag.Diagnostics
	return diags
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "->> Action Update Starting...")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "You cannot update an action.",
		Detail:   "Please, destroy it and re-create through the apply command.",
	})
	return diags
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "->> Action Delete Starting...")
	var diags diag.Diagnostics
	
	d.SetId("")
	tflog.Info(ctx, "->> Action Delete Completed Successfully...")
	return diags
}