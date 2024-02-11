package handlers

import (
	"Ideanest/pkg/database/mongodb/models"
	"Ideanest/pkg/database/mongodb/repository"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type OrganizationHandler struct {
	repository repository.OrganizationRepository
}

type memberEmail struct {
	Email string `json:"user_email"`
}

func NewOrganizationHandler(repository repository.OrganizationRepository) OrganizationHandler {
	return OrganizationHandler{repository: repository}
}

func (o OrganizationHandler) Create(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var organization models.Organization

	if err := c.BindJSON(&organization); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	insertOneResult, err := o.repository.InsertOne(ctx, organization)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"organization_id": insertOneResult.InsertedID.(primitive.ObjectID).Hex()})

}

func (o OrganizationHandler) FindById(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("organization_id"))

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	organization, err := o.repository.FindById(ctx, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id":      c.Param("organization_id"),
		"name":                 organization.Name,
		"description":          organization.Description,
		"organization_members": organization.OrganizationMembers,
	})

}

func (o OrganizationHandler) FindAll(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	organizationArr, err := o.repository.FindAll(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": organizationArr})

}

func (o OrganizationHandler) Update(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("organization_id"))

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	organization := models.Organization{Id: id}

	if err := c.BindJSON(&organization); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = o.repository.Update(ctx, organization)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id": c.Param("organization_id"),
		"name":            organization.Name,
		"description":     organization.Description,
	})

}

func (o OrganizationHandler) InviteMember(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("organization_id"))

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var newMemberEmail memberEmail

	if err := c.BindJSON(&newMemberEmail); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = o.repository.InviteMember(ctx, id, newMemberEmail.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user invited"})
}

func (o OrganizationHandler) DeleteById(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("organization_id"))

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = o.repository.DeleteById(ctx, id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organization deleted"})

}
