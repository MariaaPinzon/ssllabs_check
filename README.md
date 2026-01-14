# SSL Labs TLS Security Checker

## Descripción

Este proyecto es una aplicación web que permite verificar la seguridad TLS/SSL de un dominio utilizando la API de SSL Labs. La aplicación está construida con un backend en Go y un frontend en React.

- Analiza la configuración SSL/TLS de cualquier dominio
- Permite seleccionar si se desea utilizar la caché del servidor para obtener resultados más rápidos

## Instalación y Ejecución

### Backend (Go)

1. Navega al directorio del backend:
```bash
cd backend
```

2. Instala las dependencias (si las hay):
```bash
go mod download
```

3. Ejecuta el servidor:
```bash
go run .
```

El servidor backend se ejecutará en `http://localhost:8080`

#### Configuración avanzada

**Variable de entorno ALLOWED_ORIGIN**: Puedes configurar qué origen tiene permiso para hacer peticiones al backend mediante CORS:

```bash
# Windows PowerShell
$env:ALLOWED_ORIGIN="http://localhost:3000"; go run .

# Linux/Mac
ALLOWED_ORIGIN=http://localhost:3000 go run .

```

Si no se configura, por defecto usa `http://localhost:5173` (puerto de Vite) 

### Frontend 

1. Navega al directorio del frontend:
```bash
cd frontend
```

2. Instala las dependencias:
```bash
npm install
```

3. Ejecuta el servidor de desarrollo:
```bash
npm run dev
```

El frontend estará disponible en `http://localhost:5173` 

## Uso

1. Abre tu navegador y accede a la URL del frontend
2. Ingresa el dominio que deseas analizar (por ejemplo: `github.com`)
3. Selecciona si deseas usar la caché del servidor:
   - **Con caché**: Obtiene resultados más rápidos si el dominio ya fue analizado recientemente
   - **Sin caché**: Realiza un análisis nuevo y completo
4. Haz clic en el botón de análisis
5. Espera a que se complete el análisis (puede tardar varios minutos)
6. Revisa los resultados de seguridad TLS/SSL

## API de SSL Labs

Este proyecto utiliza la API v2 de SSL Labs. Para más información sobre la API, consulta la documentación oficial:

[SSL Labs API Documentation](https://github.com/ssllabs/ssllabs-scan/blob/master/ssllabs-api-docs-v2-deprecated.md)


## Estructura del Proyecto

```
ssllabs_check/
├── backend/
│   ├── main.go          # Punto de entrada del servidor
│   ├── models.go        # Modelos de datos
│   └── ssllabs.go       # Integración con SSL Labs API
├── frontend/
│   ├── src/
│   │   ├── App.jsx      # Componente principal
│   │   ├── Results.jsx  # Componente de resultados
│   │   └── ...
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## Autora

María Pinzon
