package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		// Nama repo image; bisa override: pulumi config set findings-api:repositoryName my-repo
		repoName := cfg.Get("repositoryName")
		if repoName == "" {
			repoName = "findings-api"
		}

		repo, err := ecr.NewRepository(ctx, "findings-api-ecr", &ecr.RepositoryArgs{
			Name:               pulumi.String(repoName),
			ImageTagMutability: pulumi.String("MUTABLE"),
			ImageScanningConfiguration: &ecr.RepositoryImageScanningConfigurationArgs{
				ScanOnPush: pulumi.Bool(true),
			},
			Tags: pulumi.StringMap{
				"Project": pulumi.String("findings-api"),
				"Managed": pulumi.String("pulumi"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("ecrRepositoryName", repo.Name)
		ctx.Export("ecrRepositoryUrl", repo.RepositoryUrl)
		return nil
	})
}
