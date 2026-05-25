Taskly Backend 🚀
API REST para gestión de tareas en equipo, construida con Go. Inspirada en herramientas como Jira y Trello.

Stack tecnológico

Go — lenguaje principal
Gin — framework HTTP
GORM — ORM para base de datos
PostgreSQL — base de datos
JWT — autenticación
Bcrypt — encriptación de contraseñas


Funcionalidades

🔐 Autenticación con JWT
👥 Gestión de equipos de trabajo
👤 Roles de miembros (admin / miembro)
✅ Tareas con estados y prioridades
📋 Asignación de tareas a miembros


Estructura del proyecto
taskly-backend/
├── main.go
├── models/
│   ├── usuario.go
│   ├── equipo.go
│   ├── miembro.go
│   ├── tarea.go
│   └── migrate.go
├── controllers/
│   ├── auth.go
│   ├── equipo.go
│   ├── miembro.go
│   └── tarea.go
├── middlewares/
│   └── auth.go
└── routes/
    └── routes.go

Instalación
1. Clona el repositorio
bashgit clone https://github.com/TU_USUARIO/taskly-backend.git
cd taskly-backend
2. Instala las dependencias
bashgo mod tidy
3. Configura las variables de entorno
Crea un archivo .env en la raíz del proyecto:
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=tu_contraseña
DB_NAME=taskly_db
DB_PORT=5432
JWT_SECRET=tu_clave_secreta
4. Crea la base de datos en PostgreSQL
sqlCREATE DATABASE taskly_db;
5. Corre el servidor
bashgo run main.go
El servidor corre en http://localhost:7000

Endpoints
Autenticación
MétodoRutaDescripciónPOST/registroCrear cuentaPOST/loginIniciar sesión

Equipos
MétodoRutaDescripciónPOST/equiposCrear equipoGET/equiposMis equiposGET/equipos/:idVer equipo con miembros y tareasDELETE/equipos/:idEliminar equipo (solo dueño)

Miembros
MétodoRutaDescripciónPOST/equipos/:id/miembrosAgregar miembroGET/equipos/:id/miembrosListar miembrosDELETE/equipos/:id/miembros/:usuarioIdEliminar miembro

Tareas
MétodoRutaDescripciónPOST/equipos/:id/tareasCrear tareaGET/equipos/:id/tareasTareas del equipoGET/tareas/:idVer tareaPUT/tareas/:idEditar tareaPATCH/tareas/:id/estadoCambiar estadoDELETE/tareas/:idEliminar tarea

Todos los endpoints excepto /registro y /login requieren el header Authorization: Bearer TOKEN


Estados de una tarea
pendiente → en_progreso → completada

Roles de miembros
admin   → puede gestionar el equipo
miembro → puede ver y trabajar en tareas

Próximamente

 Frontend en React
 Notificaciones en tiempo real
 Subida de archivos adjuntos
 Filtros y búsqueda de tareas


Autor
Desarrollado por Jeferson Gomez
