// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin/airgapped-network/cmd"
	"github.com/vmware-tanzu/tanzu-framework/cmd/cli/plugin/airgapped-network/fakes"
	"github.com/vmware-tanzu/tanzu-framework/tkg/utils"
)

const tkgversion = "v1.3.0"
const tkgImageRepo = "projects.registry.vmware.com/tkg"

var _ = Describe("DownloadTkgCompatibilityImage()", func() {
	var (
		fake = &fakes.ImgPkgClientFake{}
	)

	pullImage := &cmd.PublishImagesToTarOptions{}

	JustBeforeEach(func() {
		pullImage.ImageDetails = map[string]string{}

	})

	When("tkg-compatibility image is not present in public repo", func() {
		It("should return err", func() {
			pullImage.PkgClient = fake
			err := pullImage.DownloadTkgCompatibilityImage()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("image doesn't have any tags"))
		})
	})
	When("DownloadTkgCompatibilityImage successful", func() {
		It("should return nil", func() {
			tags := []string{"v1", "v3", "v2"}
			fake.ImgpkgTagListImageReturns(tags)
			pullImage.PkgClient = fake
			err := pullImage.DownloadTkgCompatibilityImage()
			Expect(err).ToNot(HaveOccurred())
			images := len(pullImage.ImageDetails)
			Expect(images).To(Equal(1))
		})
	})
})

var _ = Describe("DownloadTkgBomAndComponentImages()", func() {
	var (
		fake = &fakes.ImgPkgClientFake{}
	)
	pullImage := &cmd.PublishImagesToTarOptions{}
	JustBeforeEach(func() {
		pullImage.ImageDetails = map[string]string{}
		pullImage.TkgImageRepo = tkgImageRepo
		pullImage.TkgVersion = tkgversion
	})

	When("Error while downloading tkg-bom", func() {
		It("should return err", func() {
			fake.ImgpkgCopytotarReturns(errors.New(""))
			pullImage.PkgClient = fake
			tkgCompatibilityRelativeImagePath, err := pullImage.DownloadTkgBomAndComponentImages()
			Expect(err).To(HaveOccurred())
			Expect(tkgCompatibilityRelativeImagePath).To(ContainSubstring(""))
			Expect(err.Error()).To(ContainSubstring("error while downloading tkg-bom"))
		})
	})
	When("DownloadTkgBomAndComponentImages successful", func() {
		It("should return nil", func() {
			err := os.MkdirAll("./tmp", os.ModePerm)
			Expect(err).ToNot(HaveOccurred())
			fake.ImgpkgCopytotarReturns(nil)
			Expect(err).ToNot(HaveOccurred())
			pullImage.PkgClient = fake
			err = utils.CopyFile("./testdata/tkg-bom-v1.3.0.yaml", "./tmp/tkg-bom-v1.3.0.yaml")
			Expect(err).ToNot(HaveOccurred())
			tkgCompatibilityRelativeImagePath, err := pullImage.DownloadTkgBomAndComponentImages()
			Expect(err).ToNot(HaveOccurred())
			Expect(tkgCompatibilityRelativeImagePath).To(Equal("tkr-compatibility"))
			images := len(pullImage.ImageDetails)
			Expect(images).To(Equal(37))
			err = utils.DeleteFile("./tmp/tkg-bom-v1.3.0.yaml")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("DownloadTkrCompatibilityImage()", func() {
	var (
		fake = &fakes.ImgPkgClientFake{}
	)
	pullImage := &cmd.PublishImagesToTarOptions{}

	JustBeforeEach(func() {
		pullImage.ImageDetails = map[string]string{}
		pullImage.TkgImageRepo = tkgImageRepo
		pullImage.TkgVersion = tkgversion

	})

	When("tkr-compatibility image is not present in public repo", func() {
		It("should return err", func() {
			pullImage.PkgClient = fake
			list, err := pullImage.DownloadTkrCompatibilityImage("tkr-compatibility")
			_ = list
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("image doesn't have any tags"))
		})
	})
	When("DownloadTkrCompatibilityImage successful", func() {
		It("should return nil", func() {
			tags := []string{"v19"}
			err := os.MkdirAll("./tmp", os.ModePerm)
			Expect(err).ToNot(HaveOccurred())
			fake.ImgpkgTagListImageReturns(tags)
			pullImage.PkgClient = fake
			err = utils.CopyFile("./testdata/tkr-compatibility.yaml", "./tmp/tkr-compatibility.yaml")
			Expect(err).ToNot(HaveOccurred())
			list, err := pullImage.DownloadTkrCompatibilityImage("tkr-compatibility")
			_ = list
			Expect(err).ToNot(HaveOccurred())
			images := len(pullImage.ImageDetails)
			Expect(images).To(Equal(1))
			err = utils.DeleteFile("./tmp/tkr-compatibility.yaml")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("DownloadTkrBomAndComponentImages()", func() {
	var (
		fake = &fakes.ImgPkgClientFake{}
	)
	pullImage := &cmd.PublishImagesToTarOptions{}

	JustBeforeEach(func() {
		pullImage.ImageDetails = map[string]string{}
		pullImage.TkgImageRepo = tkgImageRepo
		pullImage.TkgVersion = tkgversion

	})
	When("Error while downloading tkr bom", func() {
		It("should return err", func() {
			fake.ImgpkgCopytotarReturns(errors.New("error while downloading tkr bom"))
			pullImage.PkgClient = fake
			err := pullImage.DownloadTkrBomAndComponentImages("v1.20.4+vmware.1-tkg.1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error while downloading tkr bom"))
		})
	})
	When("DownloadTkrBomAndComponentImages successful", func() {
		It("should return nil", func() {
			err := os.MkdirAll("./tmp", os.ModePerm)
			Expect(err).ToNot(HaveOccurred())
			fake.ImgpkgCopytotarReturns(nil)
			pullImage.PkgClient = fake
			err = utils.CopyFile("./testdata/tkr-bom-v1.17.16+vmware.2-tkg.1.yaml", "./tmp/tkr-bom-v1.17.16+vmware.2-tkg.1.yaml")
			Expect(err).ToNot(HaveOccurred())
			err = pullImage.DownloadTkrBomAndComponentImages("v1.17.16+vmware.2-tkg.1")
			Expect(err).ToNot(HaveOccurred())
			images := len(pullImage.ImageDetails)
			Expect(images).To(Equal(10))
			err = utils.DeleteFile("./tmp/tkr-bom-v1.17.16+vmware.2-tkg.1.yaml")
			Expect(err).ToNot(HaveOccurred())
		})
	})

})
