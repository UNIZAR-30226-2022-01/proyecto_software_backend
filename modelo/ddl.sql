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
	CONSTRAINT "Partida_pkey" PRIMARY KEY (id),
	CONSTRAINT partida_check CHECK (((("esPublica" = true) AND ("passwordHash" IS NULL)) OR (("esPublica" = false) AND ("passwordHash" IS NOT NULL))))
);


CREATE TABLE backend."Usuario" (
	id serial4 NOT NULL,
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
	CONSTRAINT "Usuario_email_key" UNIQUE (email),
	CONSTRAINT "Usuario_pkey" PRIMARY KEY (id),
	CONSTRAINT usuario_un UNIQUE ("nombreUsuario"),
	CONSTRAINT "Usuario_ID_dado_fkey" FOREIGN KEY ("ID_dado") REFERENCES backend."ItemTienda"(id),
	CONSTRAINT "Usuario_ID_ficha_fkey" FOREIGN KEY ("ID_ficha") REFERENCES backend."ItemTienda"(id)
);


CREATE TABLE backend."EsAmigo" (
	"ID_usuario1" int4 NOT NULL,
	"ID_usuario2" int4 NOT NULL,
	pendiente bool NOT NULL,
	CONSTRAINT "EsAmigo_pkey" PRIMARY KEY ("ID_usuario1", "ID_usuario2"),
	CONSTRAINT "EsAmigo_ID_usuario1_fkey" FOREIGN KEY ("ID_usuario1") REFERENCES backend."Usuario"(id),
	CONSTRAINT "EsAmigo_ID_usuario2_fkey" FOREIGN KEY ("ID_usuario2") REFERENCES backend."Usuario"(id)
);


CREATE TABLE backend."Participa" (
	"ID_partida" int4 NOT NULL,
	"ID_usuario" int4 NOT NULL,
	CONSTRAINT "Participa_pkey" PRIMARY KEY ("ID_partida", "ID_usuario"),
	CONSTRAINT "Participa_ID_partida_fkey" FOREIGN KEY ("ID_partida") REFERENCES backend."Partida"(id),
	CONSTRAINT "Participa_ID_usuario_fkey" FOREIGN KEY ("ID_usuario") REFERENCES backend."Usuario"(id)
);


CREATE TABLE backend."TieneItems" (
	"ID_item" int4 NOT NULL,
	"ID_user" int4 NOT NULL,
	CONSTRAINT "TieneItems_pkey" PRIMARY KEY ("ID_item", "ID_user"),
	CONSTRAINT "TieneItems_ID_item_fkey" FOREIGN KEY ("ID_item") REFERENCES backend."ItemTienda"(id),
	CONSTRAINT "TieneItems_ID_user_fkey" FOREIGN KEY ("ID_user") REFERENCES backend."Usuario"(id)
);
