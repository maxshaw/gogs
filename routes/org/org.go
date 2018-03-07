// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package org

import (
	log "gopkg.in/clog.v1"

	"github.com/maxshaw/gogs/models"
	"github.com/maxshaw/gogs/pkg/context"
	"github.com/maxshaw/gogs/pkg/form"
	"github.com/maxshaw/gogs/pkg/setting"
)

const RouteCreate = "org/create"

func Create(c *context.Context) {
	c.Data["Title"] = c.Tr("new_org")
	c.HTML(200, RouteCreate)
}

func CreatePost(c *context.Context, f form.CreateOrg) {
	c.Data["Title"] = c.Tr("new_org")

	if c.HasError() {
		c.HTML(200, RouteCreate)
		return
	}

	org := &models.User{
		Name:     f.OrgName,
		IsActive: true,
		Type:     models.USER_TYPE_ORGANIZATION,
	}

	if err := models.CreateOrganization(org, c.User); err != nil {
		c.Data["Err_OrgName"] = true
		switch {
		case models.IsErrUserAlreadyExist(err):
			c.RenderWithErr(c.Tr("form.org_name_been_taken"), RouteCreate, &f)
		case models.IsErrNameReserved(err):
			c.RenderWithErr(c.Tr("org.form.name_reserved", err.(models.ErrNameReserved).Name), RouteCreate, &f)
		case models.IsErrNamePatternNotAllowed(err):
			c.RenderWithErr(c.Tr("org.form.name_pattern_not_allowed", err.(models.ErrNamePatternNotAllowed).Pattern), RouteCreate, &f)
		default:
			c.Handle(500, "CreateOrganization", err)
		}
		return
	}
	log.Trace("Organization created: %s", org.Name)

	c.Redirect(setting.AppSubURL + "/org/" + f.OrgName + "/dashboard")
}
