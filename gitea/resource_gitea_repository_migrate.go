package gitea

import (
	"fmt"
	"log"

	giteaapi "code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiteaRepositoryMigrate() *schema.Resource {
	return &schema.Resource{
		Create: resourceGiteaRepositoryMigrateCreate,
		Read:   resourceGiteaRepositoryMigrateRead,
		Update: resourceGiteaRepositoryMigrateUpdate,
		Delete: resourceGiteaRepositoryMigrateDelete,
		Schema: map[string]*schema.Schema{
			"mirror_clone_addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			//"clone_addr": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			"uid": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mirror": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRepositoryMigrateToState(d *schema.ResourceData, org *giteaapi.Repository) error {
	//if err := d.Set("clone_addr", org.CloneURL); err != nil {
	//	return err
	//}
	if err := d.Set("uid", org.Owner.ID); err != nil {
		return err
	}
	if err := d.Set("description", org.Description); err != nil {
		return err
	}
	if err := d.Set("mirror", org.Mirror); err != nil {
		return err
	}
	if err := d.Set("private", org.Private); err != nil {
		return err
	}
	if err := d.Set("name", org.Name); err != nil {
		return err
	}
	//if err := d.Set("owner", org.Owner.UserName); err != nil {
	//	return err
	//}
	return nil
}

func resourceGiteaRepositoryMigrateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)

	options := giteaapi.MigrateRepoOption{
		CloneAddr:   d.Get("mirror_clone_addr").(string),
		UID:         d.Get("uid").(int),
		RepoName:    d.Get("name").(string),
		Mirror:      d.Get("mirror").(bool),
		Private:     d.Get("private").(bool),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] create mirror repo %q", options.CloneAddr)

	org, err := client.MigrateRepo(options)
	if err != nil {
		return fmt.Errorf("unable to create mirror repoanization: %v", err)
	}
	log.Printf("[DEBUG]  mirror repo created: %v", org)
	d.SetId(fmt.Sprintf("%d", org.ID))
	return resourceRepositoryMigrateToState(d, org)
}

func resourceGiteaRepositoryMigrateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	repoName := d.Get("name").(string)
	repo, err := client.GetRepo(owner, repoName)
	if err != nil {
		return fmt.Errorf("unable to retrieve repository %s %s", owner, repoName)
	}
	log.Printf("[DEBUG] repository find: %v", repo)

	return resourceRepositoryMigrateToState(d, repo)
}

func resourceGiteaRepositoryMigrateUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGiteaRepositoryMigrateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*giteaapi.Client)
	owner := d.Get("owner").(string)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] delete repository: %s %s", owner, name)
	return client.DeleteRepo(owner, name)
}
