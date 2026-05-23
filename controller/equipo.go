package controller

import (
	"strconv"
	"taskly-backend/models"

	"github.com/gin-gonic/gin"
)

// CrearEquipo      → crea un equipo, el OwnerID lo saca del token
func CrearEquipo(c *gin.Context) {

	var equipo models.Equipo
	if err := c.BindJSON(&equipo); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	ownerid := c.GetUint("id")

	equipo.OwnerID = ownerid
	DB.Create(&equipo)
	c.JSON(200, equipo)

}

//ObtenerEquipos   → lista todos los equipos del usuario autenticado

func ObtenerEquipos(c *gin.Context) {

	var equipos []models.Equipo

	ownerid := c.GetUint("id")

	resultado := DB.Where("owner_id  = ?", ownerid).Find(&equipos)

	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encontraron equipos asociados"})
		return
	}

	c.JSON(200, equipos)

}

//ObtenerEquipo    → trae un equipo por ID con sus miembros y tareas

func ObtenerEquipo(c *gin.Context) {
	var equipo models.Equipo
	idEquipo := c.Param("id")

	resultado := DB.Preload("Miembros").Preload("Tareas").First(&equipo, idEquipo)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "no se encontro ningun equipo"})
		return

	}

	c.JSON(200, equipo)

}

//EliminarEquipo   → solo el dueño puede eliminar (verifica OwnerID)

func EliminarEquipo(c *gin.Context) {

	var equipo models.Equipo

	idstr := c.Param("Id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		c.JSON(400, gin.H{"error": "Id no correcto"})
		return
	}

	resultado := DB.First(&equipo, id)

	if resultado.Error != nil {
		c.JSON(400, gin.H{"error": "No se encontro el equipo"})
		return
	}

	// verificarel usuario logugo

	if equipo.OwnerID != c.GetUint("id") {
		c.JSON(404, gin.H{"error": "No tienes permisos"})
		return

	}

	DB.Delete(&equipo)
	c.JSON(200, gin.H{"message": "Se elimino el quipo correctamente"})

}
