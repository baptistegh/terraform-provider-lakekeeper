package provider

import (
	"context"
	"fmt"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperWarehouseDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperWarehouseDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperWarehouseDataSource)
}

// NewLakekeeperWarehouseDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseDataSource() datasource.DataSource {
	return &LakekeeperWarehouseDataSource{}
}

// LakekeeperWarehouseDataSource is the data source implementation.
type LakekeeperWarehouseDataSource struct {
	client *lakekeeper.Client
}

type storageProfileDataSourceWrapper struct {
	S3StorageProfile   *sdk.S3StorageProfileDataSourceModel   `tfsdk:"s3"`
	ADLSStorageProfile *sdk.ADLSStorageProfileDataSourceModel `tfsdk:"adls"`
	GCSStorageProfile  *sdk.GCSStorageProfileDataSourceModel  `tfsdk:"gcs"`
}

// lakekeeperWarehouseDataSourceModel describes the data source data model.
type lakekeeperWarehouseDataSourceModel struct {
	ID             types.String                     `tfsdk:"id"` // form: project_id:warehouse_id (internal ID)
	WarehouseID    types.String                     `tfsdk:"warehouse_id"`
	Name           types.String                     `tfsdk:"name"`
	ProjectID      types.String                     `tfsdk:"project_id"`
	Protected      types.Bool                       `tfsdk:"protected"`
	Active         types.Bool                       `tfsdk:"active"`
	ManagedAccess  types.Bool                       `tfsdk:"managed_access"`
	StorageProfile *storageProfileDataSourceWrapper `tfsdk:"storage_profile"`
	DeleteProfile  *sdk.DeleteProfileModel          `tfsdk:"delete_profile"`
}

// Metadata returns the data source type name.
func (d *LakekeeperWarehouseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

// Schema defines the schema for the data source.
func (d *LakekeeperWarehouseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse`" + ` data source retrieves information about a warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID the warehouse. In the form: {{project_id}}/{{warehouse_id}}",
				Computed:            true,
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID the warehouse.",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID to which the warehouse belongs. If not provided, the default project will be used.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the warehouse.",
				Computed:            true,
			},
			"protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is protected from being deleted.",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is active.",
				Computed:            true,
			},
			"managed_access": schema.BoolAttribute{
				MarkdownDescription: "Whether managed access is active for this warehouse.",
				Computed:            true,
			},
			"storage_profile": schema.SingleNestedAttribute{
				MarkdownDescription: "The storage profile of the warehouse. One of `s3`, `adls` or `gcs`.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"s3": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: `S3 storage profile. Suitable for AWS or any service compatible with the S3 API.`,
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Region to use for S3 requests.",
							},
							"bucket": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The bucket name for the storage profile.",
							},
							"sts_enabled": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.",
							},
							"key_prefix": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
							"allow_alternative_protocols": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
							},
							"assume_role_arn": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile",
							},
							"aws_kms_key_arn": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "ARN of the KMS key used to encrypt the S3 bucket, if any.",
							},
							"endpoint": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`",
							},
							"flavor": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "S3 flavor to use. Defaults to `aws`.",
							},
							"path_style_access": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.",
							},
							"push_s3_delete_disabled": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Controls whether the `s3.delete-enabled=false` flag is sent to clients.",
							},
							"remote_signing_url_style": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "S3 URL style detection mode for remote signing. One of `auto`, `path-style`, `virtual-host`. Default: `auto`.",
							}, "sts_role_arn": schema.StringAttribute{
								Computed: true,
							},
							"sts_token_validity_seconds": schema.Int64Attribute{
								Computed:            true,
								MarkdownDescription: "The validity of the STS tokens in seconds. Default is `3600`.",
							},
						},
					},
					"adls": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "ADLS storage profile. Suitable for Azure Data Lake Storage Gen2.",
						Attributes: map[string]schema.Attribute{
							"account_name": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Name of the azure storage account.",
							},
							"filesystem": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Name of the adls filesystem, in blobstorage also known as container.",
							},
							"allow_alternative_protocols": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Allow alternative protocols such as wasbs:// in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
							},
							"authority_host": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The authority host to use for authentication. Defaults to `https://login.microsoftonline.com`.",
							},
							"host": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The host to use for the storage account. Defaults to `dfs.core.windows.net`.",
							},
							"key_prefix": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
							"sas_token_validity_seconds": schema.Int64Attribute{
								Computed:            true,
								MarkdownDescription: "The validity of the sas token in seconds. Default is `3600`.",
							},
						},
					},
					"gcs": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "GCS storage profile. Designed for use with Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"bucket": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The bucket name.",
							},
							"key_prefix": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
						},
					},
				},
			},
			"delete_profile": sdk.DeleteProfileDatasourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperWarehouseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperWarehouseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperWarehouseDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	warehouseID := state.WarehouseID.ValueString()
	projectID := state.ProjectID.ValueString()

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", projectID, warehouseID))

	warehouse, _, err := d.client.WarehouseV1(projectID).Get(ctx, warehouseID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s, %v", state.Name.ValueString(), err))
		return
	}

	// Authorization Properties
	m, _, err := d.client.PermissionV1().WarehousePermission().GetAuthzProperties(ctx, warehouseID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read authorization properties for warehouse %s, %v", state.Name.ValueString(), err))
		return
	}
	state.ManagedAccess = types.BoolValue(m.ManagedAccess)

	diags := state.RefreshDataSourceFromSettings(warehouse)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
