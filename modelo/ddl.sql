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


INSERT INTO backend."Usuario" (email,"nombreUsuario","passwordHash",biografia,"cookieSesion","partidasGanadas","partidasTotales",puntos,"ID_dado","ID_ficha") VALUES
	 ('creadorP4@creadorP4.com','creadorP4','$2a$10$P3pTCoGg/GkHeQkhl9CLkO6HhW.cdjZlYg2Kaz6VHMxzsXWCU84MW','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250347C4447445332336D5350594B37766D75623258387558497536464E634A6A41525154325256685672743133503669357843724C35466333476375484330336B61756E4154555052486A476D75566D784873797A7A4262794F6F6E71565553444B30563666464D547550436457703951667557576967336535675848373272384E45030F010000000ED9EBBC9A236F13E2007800','hex'),0,0,0,0,0),
	 ('creadorP5@creadorP5.com','creadorP5','$2a$10$wVtJnVfKRxLTIzldGvY.ve48W25MMB/PbkxrXiEoxGl1YCPzb2xQy','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250357C6C436E717347775671795465437464526C463238684C5A345A533655394E6C5541395771384A35384F4E654462674762563752656A73474E68693345366F6C496439494C32476878377353474748676D4B77554D6F61355334565A6C7239625546326C43437A46633438627A6C6D5733534E39445841305953794E42694F5734030F010000000ED9EBBC9A26B7EC2B007800','hex'),0,0,0,0,0),
	 ('creadorP1@creadorP1.com','creadorP1','$2a$10$E9IEVX2VfXzOvqVkgggvmO5rvcDIBAK1IUuoBF.r6kgj/DBKampMm','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250317C42704C6E6667447363325744384632714E66484B356138346A6A4A6B777A446B6839683266686655567553396A5A387556626856337643354157583339495655575350324E634863695776715A5461324E3935527852545A485755736144364845647A305468625866513670595351336E3236376C3156514B474E6253754A45030F010000000ED9EBBC9A15EC635B007800','hex'),0,0,0,0,0),
	 ('creadorP2@creadorP2.com','creadorP2','$2a$10$3rU9ga9a6PyZ0qSZmj8xC.Rkli8YjSeITzeKv6lffose3v9pENpJe','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250327C396651627A4F4E4A4141776443786D4D38424961624B4552735568504E6D4D6D64663265534A795974717763466955494C7A58763266634E4972574F3773546F46676F696C4130553157784E6557316764675556447345574A3737615837744C464A38347159553655724E3863746563775A743553347A6A684430745852546D030F010000000ED9EBBC9A1BB2D202007800','hex'),0,0,0,0,0),
	 ('creadorP3@creadorP3.com','creadorP3','$2a$10$ZirWpg7LgFG/e0JhfX68i.eYoLwv5WT9VfsqcnEd.uMM8mNxHL9q6','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250337C6B594B516F4E3931466D576E51534B327752433555484B324B714174786A50325A6D44316A7474337A6772354D65556A6F416A634F39617A4D6D7455335974763050374F506D6D534E6138376436747374617879356E61636E4A42537546704F687949584536504A304468556B4C587159596E454E756E716473777054773455030F010000000ED9EBBC9A1FD5182C007800','hex'),0,0,0,0,0),
	 ('creadorP6@creadorP6.com','creadorP6','$2a$10$LZmxjQkYbPl8pFlONxW30ul1JTk7m5B3yUfEbSG5jm.GHQjm6Egtm','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFAEFF82010B636F6F6B69655F7573657201FF8A63726561646F7250367C713452734D5164384A5773456451514E564972756473594C44415248456535675246514B7575494251313131726874734E50486344314C5A654B6C4E7865435A475242504C634F4C5172556F4748656E6D6A76334C7976664E6B6E646456324876634159654D37634F49726E6934674A707534686665414D794A5478434C3067030F010000000ED9EBBC9A2A07B37D007800','hex'),0,0,0,0,0),
	 ('userPrincipal@userPrincipal.com','userPrincipal','$2a$10$2XeF3YlB0o8ozy4J4RGxMO337Dq4yYJ1u34H2VgIQxzQeNpF8sSLq','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFB2FF82010B636F6F6B69655F7573657201FF8E757365725072696E636970616C7C794D5A305073316F786A48416139476F3338566677756B3364616351766C464A59784C6478387751516D75614859546F596A537149447746334A4D31417552394B545A3630346F6F3456683556396F47767832704773506F31594A466D6E5372396C6436716F5A4232523644335A31797A36573637647064747672774831644D030F010000000ED9EBBC9A2D617828007800','hex'),0,0,0,0,0),
	 ('amigo1@amigo1.com','amigo1','$2a$10$wuWqFkcos7STNC5MFPfvOu/usYRIPP8xsQUdhkaJdsjbsewzn9ZA6','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFABFF82010B636F6F6B69655F7573657201FF87616D69676F317C676C45685867756F353772714631596174793475714D674E6C31486A7175394772516E34477162414C4178534258776D6565334B587475563966486A7847557A526D43393867417A736A7477706B76704D526B6437597153346246434C427346544F326F57576538336E66476E7855537575766C767A47424E5269794F594930030F010000000ED9EBBC9A30B45E48007800','hex'),0,0,0,0,0),
	 ('amigo2@amigo2.com','amigo2','$2a$10$afNs6UJgHtQwCH0MOwWx/O8smzvoFQEEoh7e0fX85rv4TJqNUqWdW','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFABFF82010B636F6F6B69655F7573657201FF87616D69676F327C7A376463327676714662455A7358636952746631556C736D545937414C3975524F4C5379573143506679576D63334233386F574C43626C493542376344674F446A344F7539326A5948515A4C47444E4257545A5461543939565A4E7731794541374A4E6B563164344E4B30615771353155316A6B4764674E6959433566634A71030F010000000ED9EBBC9A34057A09007800','hex'),0,0,0,0,0),
	 ('amigo3@amigo3.com','amigo3','$2a$10$C6M5PtnU9G0pOWHyD/XCqO5ziyp5sFp89TXI5Ayq0kXDEWQUFJA7m','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFABFF82010B636F6F6B69655F7573657201FF87616D69676F337C796D6D5647333932517657734165325748576F4D7A7254437859436F4350785A3931615830763864363379635077464A33617441334139346C794C4D53444732667146753771776E5461465477684143634F4F714158356856354A423862644B6975615755615830446B46496972636465393836454B78477A4B544A4C6A4343030F010000000ED9EBBC9A3850E349007800','hex'),0,0,0,0,0);
INSERT INTO backend."Usuario" (email,"nombreUsuario","passwordHash",biografia,"cookieSesion","partidasGanadas","partidasTotales",puntos,"ID_dado","ID_ficha") VALUES
	 ('amigo4@amigo4.com','amigo4','$2a$10$J43ry3aZ8ujskW3cWTG2G.be..bmuPgsBtJEdO2WFWNtGgophYcue','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFABFF82010B636F6F6B69655F7573657201FF87616D69676F347C6E3934514533306E7345477A4635676D4636394531426C413050503134544C4959683752317A56394443644536506E4531384573566736675A717A5938327958337762316C786359676B4F6535494E34586562537A794446776F65706E784B62654C4946623236766974714D76764C7779544D3744345566354F706D6D725853030F010000000ED9EBBC9B000D937F007800','hex'),0,0,0,0,0),
	 ('amigo5@amigo5.com','amigo5','$2a$10$FHriVUusaNY9Z1D7D.ub2uTz2Brb3DUVwIB8JF4zuyVzJG8ZmDUBa','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFABFF82010B636F6F6B69655F7573657201FF87616D69676F357C6D7535476833555758415A38555772777A6743556E30773342726C6E4B68644447785363484C474D6E43536C6C526B546D32684866454B595A45784970784E45563635484D515A324C6830444149323659306D41786C754E473975634B7852454D45577A7654635A66636B647255617839634D7078304F365262745631685875030F010000000ED9EBBC9B040A11C7007800','hex'),0,0,0,0,0),
	 ('NoAmigo1@NoAmigo1.com','NoAmigo1','$2a$10$eYuP5Eb5/H1yYoWMioUPWew.12oA.BA2KLayofE6IaNlmJoxaD.Iy','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F317C557031446E4561373656726B7636764646664851375A4E5871584F31346E7559773534384C715A4449693934616F376F6F386B334469664E3836344F336F6B356A546D716358665A7332584D6B524B51667049776F535951446E4972324377635A5A663264554D3948306576626E4470766A41714F74634F6A6553314C4E6471030F010000000ED9EBBC9B094B16EF007800','hex'),0,0,0,0,0),
	 ('NoAmigo2@NoAmigo2.com','NoAmigo2','$2a$10$SXk/cVO8ltn/.3viWAz.L.Y.wld3vTaACvvINjwGMkMehBVj.WTUy','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F327C3766677759344A5A42513172644F6A563639357647737A4E6D5642716378714171657A414C75503148703250506E48596C75655257474B71435A57484B4F6869715861307768796C7751456A586A6F77753466424B5375365063623346375A6F6A394D51727664393873735565436B707159686A52646A784D486F4D524F4B37030F010000000ED9EBBC9B0D08657F007800','hex'),0,0,0,0,0),
	 ('NoAmigo3@NoAmigo3.com','NoAmigo3','$2a$10$c.fCwfBEn0uIXs36.uPhh.Zx5ct3Du0L34ymNLGymMWZYNHETIuX6','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F337C703930434B77324A306562396977524F73635055346B7A6A754679384F3257425738414437736378595648383773517742616D765746794C78684B3834645967476D56323472384A594444613436716D326650505941744379624378544C69594F6135637164336669394B636D374F31375A7534736F6A424E6D5138664E5573030F010000000ED9EBBC9B124D0D9D007800','hex'),0,0,0,0,0),
	 ('NoAmigo4@NoAmigo4.com','NoAmigo4','$2a$10$cHdp5VelDpsupstB3UFujum1gK3GfPIpPrw5vf0geacuXBGdf6Yn2','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F347C75715A714E417244676C4A4F455969346E7A773144556177655946586A427A64693759685449674656595877595070634471575757653061674D4364616E704C78525866485231524A514269336B61473132564354686B735A4D304F7A4B3142516E7A6D74734669774538545156746D764B5258484D6D46716844313069344F030F010000000ED9EBBC9B16BC1508007800','hex'),0,0,0,0,0),
	 ('NoAmigo5@NoAmigo5.com','NoAmigo5','$2a$10$MrNDlsbE6DmDNCo9tJaAHulHOM3dm9DXh9DCiZZMHv2XMrC/oZTsS','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F357C4B37314F556E4E76587352797532765558754533366A6B53626C4C6F755559414A6E4A7242526139706936594845594A733778564F72354A7867576675505771316E4146645773424E466A5A53324156686A436854486C4F49676D7569506F6243554A306B58776F6E4F6D577466684E326D52313134316B33565A4671465462030F010000000ED9EBBC9B1A229248007800','hex'),0,0,0,0,0),
	 ('NoAmigo6@NoAmigo6.com','NoAmigo6','$2a$10$fZ5jYGl8a./Qb2ok5imt6.NBY83QEEO6v3uEVzRdvJnv/f34HoHGK','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F367C563257346C6951524B4D66516D6431386F446B737954637041625943334E52714E414F7A504749737030424F46556945734F726158476D4A5275674F5A6C526631476930467A6378316B37457245434954573843784B763956704864573679354B3570487873354D384F333371497556794A4833445136503159566167366B4B030F010000000ED9EBBC9B1FB7A5DE007800','hex'),0,0,0,0,0),
	 ('NoAmigo7@NoAmigo7.com','NoAmigo7','$2a$10$dic6z0szp1GwsJh2scxYWOaTwE0yUnOtmXkgWxms1iFWd6EEAlNHS','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F377C6F536173686C4555536958696579346438504F4D697445646B4B5A706441614B49596A30387475487A457479794A6569734F4E6B61554E50627163415675533649324F5065536F32314E315A4744326266367568476E6A34363067304471304F54354D61454C6561514A667351587548626843426148374F544A4D594F674634030F010000000ED9EBBC9B2589415A007800','hex'),0,0,0,0,0),
	 ('NoAmigo8@NoAmigo8.com','NoAmigo8','$2a$10$1.RMvACuwm6IZgmDyez.TOiGxy2JlCwo0q/fEWWEniXL7gCRBiInS','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF894E6F416D69676F387C4B714E4F6B32704E347953575835326B61356D497848546457694A38543968796C3876634975733052346353786B304B767265445836383775506E6A704E4755446D524965524C76487276724C426150336B427A644F5A504946696168625556584A634D304F63436A7332633032304C7436575173764E394243436274514C41030F010000000ED9EBBC9B2B03AD7F007800','hex'),0,0,0,0,0);

INSERT INTO backend."Partida" ("estadoPartida",mensajes,"esPublica","passwordHash","enCurso","maxJugadores") VALUES
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),false,'$2a$10$7Ob4behNR0GebHMEW/LRx.ne5xijDxXTxFs9AWDzlyvTRPZ9jU2iK',false,6),
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),false,'$2a$10$HlylQJdo2rEjZ8IuDOhDNeNBHqjn1JE9f8t8qFthLGgZSO.JfjdXu',false,6),
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),false,'$2a$10$HS/UWzV4TkSG6D2OXaTpFu2/Vbrg/IWVyPu0eCBzk1KRPGlgiLINO',false,6),
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),true,NULL,false,6),
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),true,NULL,false,6),
	 (decode('FFD8FF870301010D45737461646F5061727469646101FF8800010C0108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F72010400010446617365010400010B4E756D65726F5475726E6F010400010A45737461646F4D61706101FF9600010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F436172746101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000046FF95040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9600010401FF94000027FF93030102FF9400010201084F637570616E7465010C0001094E756D54726F706173010400000003FF8800','hex'),decode('0DFF99020102FF9A0001FF98000013FF97030101074D656E73616A6501FF9800000004FF9A0000','hex'),true,NULL,false,6);

INSERT INTO backend."Participa" ("ID_partida","nombreUsuario") VALUES
	 (1,'creadorP1'),
	 (2,'creadorP2'),
	 (3,'creadorP3'),
	 (4,'creadorP4'),
	 (5,'creadorP5'),
	 (6,'creadorP6'),
	 (1,'amigo1'),
	 (1,'amigo2'),
	 (1,'NoAmigo1'),
	 (2,'amigo3');

INSERT INTO backend."Participa" ("ID_partida","nombreUsuario") VALUES
	 (2,'NoAmigo2'),
	 (2,'NoAmigo3'),
	 (3,'NoAmigo4'),
	 (4,'amigo4'),
	 (4,'amigo5'),
	 (4,'NoAmigo5'),
	 (5,'NoAmigo6'),
	 (5,'NoAmigo7'),
	 (6,'NoAmigo8');

INSERT INTO backend."EsAmigo" ("nombreUsuario1","nombreUsuario2",pendiente) VALUES
	 ('userPrincipal','amigo1',false),
	 ('userPrincipal','amigo2',false),
	 ('userPrincipal','amigo3',false),
	 ('userPrincipal','amigo4',false),
	 ('userPrincipal','amigo5',false);
