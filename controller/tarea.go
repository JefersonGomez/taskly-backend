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
	tarea.Estado = "backlog"
	DB.Create(&tarea)
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
	DB.Where("equipo_id = ?", idEquipo).Find(&tareas)

	// para cada tarea buscar el usuario asignado
	type TareaConAsignado struct {
		models.Tarea
		NombreAsignado string `json:"nombreAsignado"`
	}

	var result []TareaConAsignado
	for _, t := range tareas {
		item := TareaConAsignado{Tarea: t}
		if t.AsignadoID != 0 {
			var usuario models.Usuario
			DB.First(&usuario, t.AsignadoID)
			item.NombreAsignado = usuario.Nombre
		}
		result = append(result, item)
	}

	c.JSON(200, result)
}

func ObtenerTarea(c *gin.Context) {
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
	if body.Estado != "backlog" && body.Estado != "pendiente" && body.Estado != "en_progreso" && body.Estado != "completada" {
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
