package controller

import (
	"strconv"
	"taskly-backend/models"

	"github.com/gin-gonic/gin"
)

func CrearTarea(c *gin.Context) {
	var equipo models.Equipo
	id := c.Param("id")
	idEquipo, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de equipo incorrecto"})
		return
	}
	resultado := DB.First(&equipo, idEquipo)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var tarea models.Tarea
	if err := c.BindJSON(&tarea); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	tarea.EquipoID = uint(idEquipo)
	tarea.Estado = "pendiente"
	DB.Create(&tarea)

	// retornar la tarea con el asignado cargado
	DB.Preload("Asignado").First(&tarea, tarea.ID)
	c.JSON(201, tarea)
}

func ObtenerTareas(c *gin.Context) {
	id := c.Param("id")
	idEquipo, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de equipo incorrecto"})
		return
	}
	var equipo models.Equipo
	resultado := DB.First(&equipo, idEquipo)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var tareas []models.Tarea
	DB.Preload("Asignado").Where("equipo_id = ?", idEquipo).Find(&tareas)
	c.JSON(200, tareas)
}

func ObtenerTarea(c *gin.Context) {
	id := c.Param("id")
	idTarea, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de tarea incorrecto"})
		return
	}
	var tarea models.Tarea
	resultado := DB.Preload("Asignado").First(&tarea, idTarea)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "tarea no encontrada"})
		return
	}
	c.JSON(200, tarea)
}

func EditarTarea(c *gin.Context) {
	id := c.Param("id")
	idTarea, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de tarea incorrecto"})
		return
	}
	var tarea models.Tarea
	resultado := DB.First(&tarea, idTarea)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "tarea no encontrada"})
		return
	}
	if err := c.BindJSON(&tarea); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	DB.Save(&tarea)
	DB.Preload("Asignado").First(&tarea, tarea.ID)
	c.JSON(200, tarea)
}

func CambiarEstado(c *gin.Context) {
	id := c.Param("id")
	idTarea, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de tarea incorrecto"})
		return
	}
	var tarea models.Tarea
	resultado := DB.First(&tarea, idTarea)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "tarea no encontrada"})
		return
	}
	var body struct {
		Estado string `json:"estado"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	if body.Estado != "pendiente" && body.Estado != "en_progreso" && body.Estado != "completada" {
		c.JSON(400, gin.H{"error": "estado no permitido"})
		return
	}
	DB.Model(&tarea).Update("estado", body.Estado)
	c.JSON(200, gin.H{"message": "estado actualizado correctamente"})
}

func EliminarTarea(c *gin.Context) {
	id := c.Param("id")
	idTarea, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "id de tarea incorrecto"})
		return
	}
	var tarea models.Tarea
	resultado := DB.First(&tarea, idTarea)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "tarea no encontrada"})
		return
	}
	DB.Delete(&tarea)
	c.JSON(200, gin.H{"message": "tarea eliminada correctamente"})
}
