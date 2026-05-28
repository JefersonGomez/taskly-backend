package routes

import (
	"taskly-backend/controller"
	"taskly-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(c *gin.Engine) {

	// 🔹 Rutas públicas (sin token)
	c.POST("/registro", controller.Registro)
	c.POST("/login", controller.Login)

	// 🔹 Endpoint de prueba para emails (sin token, solo desarrollo)
	c.POST("/debug/test-email", controller.TestEmail)

	// 🔹 Rutas protegidas (requieren token)
	p := c.Group("/")
	p.Use(middlewares.ValidarToken)
	{
		// Equipos
		p.POST("/equipos", controller.CrearEquipo)
		p.GET("/equipos", controller.ObtenerEquipos)
		p.GET("/equipos/:id", controller.ObtenerEquipo)
		p.DELETE("/equipos/:id", controller.EliminarEquipo)

		// Miembros
		p.POST("/equipos/:id/miembros", controller.AgregarMiembro)
		p.GET("/equipos/:id/miembros", controller.ObtenerMiembros)
		p.DELETE("/equipos/:id/miembros/:usuarioId", controller.EliminarMiembro)

		// Tareas
		p.POST("/equipos/:id/tareas", controller.CrearTarea)
		p.GET("/equipos/:id/tareas", controller.ObtenerTareas)
		p.GET("/tareas/:id", controller.ObtenerTarea)
		p.PUT("/tareas/:id", controller.EditarTarea)
		p.PATCH("/tareas/:id/estado", controller.CambiarEstado)
		p.DELETE("/tareas/:id", controller.EliminarTarea)

		// Usuarios
		p.GET("/usuarios/buscar", controller.BuscarUsuario)
		p.GET("/me", controller.ObtenerPerfil)
		p.PUT("/me", controller.ActualizarPerfil)
	}
}
