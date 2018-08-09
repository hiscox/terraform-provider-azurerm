package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmEaSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEaSubscriptionCreate,
		Read:   resourceArmEaSubscriptionRead,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enrollment_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"offer_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subscription_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owners": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmEaSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subscriptionClient
	ctx := meta.(*ArmClient).StopContext

	displayName := d.Get("display_name").(string)
	enrollmentAccountName := d.Get("enrollment_account_name").(string)
	offerType := subscription.OfferType(d.Get("offer_type").(string))

	creationParameters := subscription.CreationParameters{
		DisplayName: &displayName,
		OfferType:   offerType,
	}

	future, err := client.CreateSubscriptionInEnrollmentAccount(ctx, enrollmentAccountName, creationParameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	creationResult, err := future.Result(client)
	if err != nil {
		return err
	}
	d.Set("subscription_link", *creationResult.SubscriptionLink)

	return resourceArmEaSubscriptionRead(d, meta)
}

func resourceArmEaSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}