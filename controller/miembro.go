package controller

import (
	"strconv"
	"taskly-backend/models"

	"github.com/gin-gonic/gin"
)

func AgregarMiembro(c *gin.Context) {
	var equipo models.Equipo
	idEquipo := c.Param("id")
	id, err := strconv.Atoi(idEquipo)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var miembro models.Miembro
	if err := c.BindJSON(&miembro); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	if miembro.Rol != "admin" && miembro.Rol != "miembro" {
		c.JSON(400, gin.H{"error": "rol inválido, debe ser admin o miembro"})
		return
	}
	miembro.EquipoID = equipo.ID
	DB.Create(&miembro)

	// retornar el miembro con el usuario cargado
	DB.Preload("Usuario").First(&miembro, miembro.ID)
	c.JSON(201, miembro)
}

func ObtenerMiembros(c *gin.Context) {
	idEquipo := c.Param("id")
	id, err := strconv.Atoi(idEquipo)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}
	var equipo models.Equipo
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var miembros []models.Miembro
	DB.Preload("Usuario").Where("equipo_id = ?", id).Find(&miembros)
	c.JSON(200, miembros)
}

func EliminarMiembro(c *gin.Context) {
	var miembro models.Miembro
	idEquipo := c.Param("id")
	idUsuario := c.Param("usuarioId")
	id, err := strconv.Atoi(idEquipo)
	idUsu, err2 := strconv.Atoi(idUsuario)
	if err != nil || err2 != nil {
		c.JSON(400, gin.H{"error": "ids incorrectos"})
		return
	}
	resultado := DB.Where("equipo_id = ? AND usuario_id = ?", id, idUsu).First(&miembro)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "miembro no encontrado en ese equipo"})
		return
	}
	DB.Delete(&miembro)
	c.JSON(200, gin.H{"message": "miembro eliminado correctamente"})
}
