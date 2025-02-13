# Skript zur Erstellung der Verzeichnisstruktur für XMBig

# Hauptverzeichnis
$root = "."

# Vue.js Frontend
$src = "$root/src"
$assets = "$src/assets"
$components = "$src/components"
$router = "$src/router"
$store = "$src/store"
$utils = "$src/utils"

# Go Backend
$go = "$root/go"

# Protobuf
$xmbig = "$root/xmbig"

# Erstelle die Verzeichnisse
New-Item -ItemType Directory -Force -Path $src
New-Item -ItemType Directory -Force -Path $assets
New-Item -ItemType Directory -Force -Path $components
New-Item -ItemType Directory -Force -Path $router
New-Item -ItemType Directory -Force -Path $store
New-Item -ItemType Directory -Force -Path $utils
New-Item -ItemType Directory -Force -Path $go
New-Item -ItemType Directory -Force -Path $xmbig

# Erstelle leere Dateien (optional)
New-Item -ItemType File -Force -Path "$src/App.vue"
New-Item -ItemType File -Force -Path "$src/main.js"
New-Item -ItemType File -Force -Path "$router/index.js"
New-Item -ItemType File -Force -Path "$store/index.js"
New-Item -ItemType File -Force -Path "$utils/clientStatusFormatters.js"
New-Item -ItemType File -Force -Path "$go/GoAPI.go"
New-Item -ItemType File -Force -Path "$go/ClientConnector.go"
New-Item -ItemType File -Force -Path "$go/go.mod"
New-Item -ItemType File -Force -Path "$xmbig/miner.pb.go"
New-Item -ItemType File -Force -Path "$root/miner.proto"
New-Item -ItemType File -Force -Path "$root/package.json"
New-Item -ItemType File -Force -Path "$root/package-lock.json"
New-Item -ItemType File -Force -Path "$root/README.md"

Write-Host "Verzeichnisstruktur für XMBig wurde erstellt."