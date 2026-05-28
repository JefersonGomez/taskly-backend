package controller

import (
	"strconv"
	"taskly-backend/models"

	"github.com/gin-gonic/gin"
)

// CrearTarea godoc
// @Summary      Crear una nueva tarea en un equipo
// @Description  Crea una tarea dentro de un equipo específico y la inicializa en estado "backlog". Opcionalmente se puede asignar a un miembro.
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path           int             true  "ID del equipo"              example(1)
// @Param        datos  body           models.Tarea    true  "Datos de la tarea a crear"  example({"titulo": "Fix bug", "descripcion": "Corregir error en login", "asignado_id": 5})
// @Success      201    {object}       models.Tarea    "Tarea creada exitosamente"
// @Failure      400    {object}       models.ErrorResponse  "ID de equipo inválido o datos incorrectos"
// @Failure      401    {object}       models.ErrorResponse  "No autenticado"
// @Failure      404    {object}       models.ErrorResponse  "Equipo no encontrado"
// @Failure      500    {object}       models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id}/tareas [post]
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

// ObtenerTareas godoc
// @Summary      Listar tareas de un equipo
// @Description  Retorna todas las tareas asociadas a un equipo específico, incluyendo el nombre del usuario asignado (si existe).
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID del equipo"  example(1)
// @Success      200  {array}  models.TareaConAsignadoResponse  "Lista de tareas obtenida exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID de equipo inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Equipo no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id}/tareas [get]
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

// ObtenerTarea godoc
// @Summary      Obtener una tarea por ID
// @Description  Retorna los detalles completos de una tarea específica identificada por su ID.
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID de la tarea"  example(42)
// @Success      200  {object}  models.Tarea  "Tarea obtenida exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID de tarea inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Tarea no encontrada"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /tareas/{id} [get]
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

// EditarTarea godoc
// @Summary      Actualizar una tarea existente
// @Description  Permite modificar los campos de una tarea (título, descripción, asignado, etc.). Solo campos enviados serán actualizados.
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path           int             true  "ID de la tarea"     example(42)
// @Param        datos  body           models.Tarea    true  "Campos a actualizar"  example({"titulo": "Nuevo título", "asignado_id": 3})
// @Success      200    {object}       models.Tarea    "Tarea actualizada exitosamente"
// @Failure      400    {object}       models.ErrorResponse  "ID inválido o datos incorrectos"
// @Failure      401    {object}       models.ErrorResponse  "No autenticado"
// @Failure      404    {object}       models.ErrorResponse  "Tarea no encontrada"
// @Failure      500    {object}       models.ErrorResponse  "Error interno del servidor"
// @Router       /tareas/{id} [put]
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

// CambiarEstado godoc
// @Summary      Cambiar el estado de una tarea
// @Description  Actualiza únicamente el campo `estado` de una tarea. Valores permitidos: "backlog", "pendiente", "en_progreso", "completada".
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                      true  "ID de la tarea"           example(42)
// @Param        body  body  models.UpdateEstadoRequest  true  "Nuevo estado de la tarea"
// @Success      200   {object}  models.MessageResponse  "Estado actualizado exitosamente"
// @Failure      400   {object}  models.ErrorResponse  "ID inválido, datos incorrectos o estado no permitido"
// @Failure      401   {object}  models.ErrorResponse  "No autenticado"
// @Failure      404   {object}  models.ErrorResponse  "Tarea no encontrada"
// @Failure      500   {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /tareas/{id}/estado [patch]
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

// EliminarTarea godoc
// @Summary      Eliminar una tarea
// @Description  Elimina permanentemente una tarea del sistema. Esta acción no se puede deshacer.
// @Tags         Tareas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID de la tarea a eliminar"  example(42)
// @Success      200  {object}  models.MessageResponse  "Tarea eliminada exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID de tarea inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      403  {object}  models.ErrorResponse  "Permisos insuficientes para eliminar"
// @Failure      404  {object}  models.ErrorResponse  "Tarea no encontrada"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /tareas/{id} [delete]
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
