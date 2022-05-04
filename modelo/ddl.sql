CREATE SCHEMA backend AUTHORIZATION postgres;

CREATE TYPE backend.item AS ENUM ('ficha', 'dado', 'avatar');

CREATE TABLE backend."ItemTienda" (
	id int4 NOT NULL,
	nombre varchar NOT NULL,
	descripcion varchar NOT NULL,
	precio int4 NOT NULL,
	tipo backend.item NOT NULL,
	CONSTRAINT "ItemTienda_pkey" PRIMARY KEY (id)
);


CREATE TABLE backend."Partida" (
	id serial4 NOT NULL,
	"estadoPartida" bytea NOT NULL,		-- Serializado desde golang
	mensajes bytea NOT NULL,			-- Serializado desde golang
	"esPublica" bool NOT NULL,
	"passwordHash" varchar NULL,
	"enCurso" bool NOT NULL,
	"maxJugadores" int NOT NULL,
	CONSTRAINT "Partida_pkey" PRIMARY KEY (id),
	CONSTRAINT partida_check CHECK (((("esPublica" = true) AND ("passwordHash" IS NULL)) OR (("esPublica" = false) AND ("passwordHash" IS NOT NULL))))
);


CREATE TABLE backend."Usuario" (
	email varchar NOT NULL,
	"nombreUsuario" varchar NOT NULL,
	"passwordHash" varchar NOT NULL,
	biografia varchar NULL,
	"cookieSesion" bytea NOT NULL,		-- Serializado desde golang
	"partidasGanadas" int4 NOT NULL,
	"partidasTotales" int4 NOT NULL,
	puntos int4 NOT NULL,
	"ID_dado" int4 NOT NULL,
	"ID_ficha" int4 NOT NULL,
	"notificacionesPendientesConEstado" bytea NULL, -- Serializado desde golang
	"tokenResetPassword" varchar NULL,
	"ultimaPeticionResetPassword" date NULL,
	CONSTRAINT "Usuario_email_key" UNIQUE (email),
	CONSTRAINT "Usuario_pkey" PRIMARY KEY ("nombreUsuario"),
	CONSTRAINT usuario_un UNIQUE ("nombreUsuario"),
	CONSTRAINT "Usuario_ID_dado_fkey" FOREIGN KEY ("ID_dado") REFERENCES backend."ItemTienda"(id),
	CONSTRAINT "Usuario_ID_ficha_fkey" FOREIGN KEY ("ID_ficha") REFERENCES backend."ItemTienda"(id)
);


CREATE TABLE backend."EsAmigo" (
	"nombreUsuario1" varchar NOT NULL,
	"nombreUsuario2" varchar NOT NULL,
	pendiente bool NOT NULL,
	CONSTRAINT "EsAmigo_pkey" PRIMARY KEY ("nombreUsuario1", "nombreUsuario2"),
	CONSTRAINT "EsAmigo_nombre_usuario1_fkey" FOREIGN KEY ("nombreUsuario1") REFERENCES backend."Usuario"("nombreUsuario"),
	CONSTRAINT "EsAmigo_nombre,usuario2_fkey" FOREIGN KEY ("nombreUsuario2") REFERENCES backend."Usuario"("nombreUsuario")
);


CREATE TABLE backend."Participa" (
	"ID_partida" int4 NOT NULL,
	"nombreUsuario" varchar NOT NULL,
	CONSTRAINT "Participa_pkey" PRIMARY KEY ("ID_partida", "nombreUsuario"),
	CONSTRAINT "Participa_ID_partida_fkey" FOREIGN KEY ("ID_partida") REFERENCES backend."Partida"(id) ON DELETE CASCADE ,
	CONSTRAINT "Participa_ID_usuario_fkey" FOREIGN KEY ("nombreUsuario") REFERENCES backend."Usuario"("nombreUsuario") ON DELETE CASCADE
);


CREATE TABLE backend."TieneItems" (
	"ID_item" int4 NOT NULL,
	"nombreUsuario" varchar NOT NULL,
	CONSTRAINT "TieneItems_pkey" PRIMARY KEY ("ID_item", "nombreUsuario"),
	CONSTRAINT "TieneItems_ID_item_fkey" FOREIGN KEY ("ID_item") REFERENCES backend."ItemTienda"(id),
	CONSTRAINT "TieneItems_ID_user_fkey" FOREIGN KEY ("nombreUsuario") REFERENCES backend."Usuario"("nombreUsuario")
);


INSERT INTO backend."ItemTienda"
(id, nombre, descripcion, precio, tipo)
VALUES(0, 'Fichas por defecto.', 'Fichas por defecto, aburridas.', 0, 'ficha'::backend."item");
INSERT INTO backend."ItemTienda"
(id, nombre, descripcion, precio, tipo)
VALUES(1, 'Dados por defecto.', 'Dados por defecto, aburridos.', 0, 'dado'::backend."item");
INSERT INTO backend."ItemTienda"
(id, nombre, descripcion, precio, tipo)
VALUES(2, 'Avatar por defecto.', 'Avatar por defecto, aburrido.', 0, 'avatar'::backend."item");

INSERT INTO backend."ItemTienda" (id, nombre, descripcion, precio, tipo) VALUES
     (4, 'Fichas rojas', 'Fichas de color rojo', 10, 'ficha'),
     (5, 'Fichas azules', 'Fichas de color azul', 10, 'ficha'),
     (6, 'Fichas verdes', 'Fichas de color rojo', 10, 'ficha'),
     (7, 'Dados plateados', 'Dados de color plateado', 10, 'dado'),
     (8, 'Dados dorados', 'Dados de color dorado', 10, 'dado');