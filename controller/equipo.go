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
	equipo.OwnerID = c.GetUint("id")
	DB.Create(&equipo)
	c.JSON(201, equipo)
}

func ObtenerEquipos(c *gin.Context) {
	usuarioID := c.GetUint("id")

	var equiposDueno []models.Equipo
	DB.Preload("Miembros.Usuario").Preload("Tareas").Where("owner_id = ?", usuarioID).Find(&equiposDueno)

	var miembros []models.Miembro
	DB.Where("usuario_id = ?", usuarioID).Find(&miembros)

	var equiposMiembro []models.Equipo
	for _, m := range miembros {
		var equipo models.Equipo
		DB.Preload("Miembros.Usuario").Preload("Tareas").First(&equipo, m.EquipoID)
		equiposMiembro = append(equiposMiembro, equipo)
	}

	equipos := append(equiposDueno, equiposMiembro...)
	c.JSON(200, equipos)
}

func ObtenerEquipo(c *gin.Context) {
	var equipo models.Equipo
	idEquipo := c.Param("id")
	resultado := DB.Preload("Miembros.Usuario").Preload("Tareas").First(&equipo, idEquipo)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	c.JSON(200, equipo)
}

func EliminarEquipo(c *gin.Context) {
	var equipo models.Equipo
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	if equipo.OwnerID != c.GetUint("id") {
		c.JSON(403, gin.H{"error": "no tienes permisos para eliminar este equipo"})
		return
	}
	DB.Delete(&equipo)
	c.JSON(200, gin.H{"message": "equipo eliminado correctamente"})
}
