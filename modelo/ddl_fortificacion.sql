-- Fase de fortificación
-- Usuarios "jugadorx" con x del 1 al 6
-- Cada uno con contraseña igual a su nombre
-- Le corresponde el turno al jugador6
INSERT INTO backend."Usuario" (email,"nombreUsuario","passwordHash",biografia,"cookieSesion","partidasGanadas","partidasTotales",puntos,"ID_avatar", "ID_dado", "notificacionesPendientesConEstado") VALUES
      ('jugador1@jugador1.com','jugador1','$2a$10$EPJQdiRQgLXz3bsdmKHbGebBdh6jQMqbzpN.vGafj00bFZU.W8R76','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72317C42704C6E6667447363325744384632714E66484B356138346A6A4A6B777A446B6839683266686655567553396A5A387556626856337643354157583339495655575350324E634863695776715A5461324E3935527852545A485755736144364845647A305468625866513670595351336E3236376C3156514B474E6253754A45030F010000000EDA09DDD7200D3A7A007800','hex'),0,0,0,1,9,NULL),
      ('jugador2@jugador2.com','jugador2','$2a$10$vmN4CnB5FNkO6MUP9.NNxeTsCpNmltVpylUZ5sl7kdRwOzXJgoeQW','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72327C396651627A4F4E4A4141776443786D4D38424961624B4552735568504E6D4D6D64663265534A795974717763466955494C7A58763266634E4972574F3773546F46676F696C4130553157784E6557316764675556447345574A3737615837744C464A38347159553655724E3863746563775A743553347A6A684430745852546D030F010000000EDA09DDD729CDE742007800','hex'),0,0,0,1,9,NULL),
      ('jugador3@jugador3.com','jugador3','$2a$10$RoGZWlR5laUczFae4mk8ReWgi3A.h4pDca0Hw74RGjDkrNguIjZD6','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72337C6B594B516F4E3931466D576E51534B327752433555484B324B714174786A50325A6D44316A7474337A6772354D65556A6F416A634F39617A4D6D7455335974763050374F506D6D534E6138376436747374617879356E61636E4A42537546704F687949584536504A304468556B4C587159596E454E756E716473777054773455030F010000000EDA09DDD732B83D41007800','hex'),0,0,0,1,9,NULL),
      ('jugador4@jugador4.com','jugador4','$2a$10$.4WM2NCLrM3CT8wfLfYYG.mJBBviHilmBYTaiyZiLVP4cw23DFAaO','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72347C4447445332336D5350594B37766D75623258387558497536464E634A6A41525154325256685672743133503669357843724C35466333476375484330336B61756E4154555052486A476D75566D784873797A7A4262794F6F6E71565553444B30563666464D547550436457703951667557576967336535675848373272384E45030F010000000EDA09DDD73781EBCA007800','hex'),0,0,0,1,9,NULL),
      ('jugador5@jugador5.com','jugador5','$2a$10$TsLasNmaAaegRsNK0oQ1Bu27/V3ELiproQMMqrIFWB/peHJv8lGkq','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72357C6C436E717347775671795465437464526C463238684C5A345A533655394E6C5541395771384A35384F4E654462674762563752656A73474E68693345366F6C496439494C32476878377353474748676D4B77554D6F61355334565A6C7239625546326C43437A46633438627A6C6D5733534E39445841305953794E42694F5734030F010000000EDA09DDD800B13602007800','hex'),0,0,0,1,9,NULL),
      ('jugador6@jugador6.com','jugador6','$2a$10$kzitg/jP1qGGX90yPxVuduLF0rmGukx9Jhkt7RoW4acbkhoG1UIPW','',decode('FF9DFF8103010106436F6F6B696501FF8200010C01044E616D65010C00010556616C7565010C00010450617468010C000106446F6D61696E010C0001074578706972657301FF8400010A52617745787069726573010C0001064D617841676501040001065365637572650102000108487474704F6E6C79010200010853616D65536974650104000103526177010C000108556E70617273656401FF8600000010FF830501010454696D6501FF8400000016FF85020101085B5D737472696E6701FF8600010C0000FFADFF82010B636F6F6B69655F7573657201FF896A756761646F72367C713452734D5164384A5773456451514E564972756473594C44415248456535675246514B7575494251313131726874734E50486344314C5A654B6C4E7865435A475242504C634F4C5172556F4748656E6D6A76334C7976664E6B6E646456324876634159654D37634F49726E6934674A707534686665414D794A5478434C3067030F010000000EDA09DDD805AA1C1D007800','hex'),0,0,0,1,9,NULL);

INSERT INTO backend."Partida" ("estadoPartida",mensajes,"esPublica","passwordHash","enCurso","maxJugadores") VALUES
	 (decode('FE01C1FF870301010D45737461646F5061727469646101FF880001160108416363696F6E657301FF8A0001094A756761646F72657301FF8600011045737461646F734A756761646F72657301FF9200010C5475726E6F4A756761646F7201040001104A756761646F72657341637469766F7301FF9400010446617365010400010A45737461646F4D61706101FF9800010643617274617301FF9000010944657363617274657301FF9000010A4E756D43616D62696F73010400010D4861436F6E717569737461646F010200010F4861526563696269646F4361727461010200010D4861466F72746966696361646F0102000112526567696F6E556C74696D6F41746171756501040001114461646F73556C74696D6F417461717565010400011A54726F7061735065726469646173556C74696D6F41746171756501040001174861795465727269746F72696F4465736F63757061646F010200010E556C74696D6F446566656E736F72010C0001095465726D696E616461010200011E4A756761646F72657352657374616E746573506F72436F6E73756C74617201FF8600010C556C74696D61416363696F6E01FF8400010D416C65727461456E766961646101020000001CFF890201010E5B5D696E74657266616365207B7D01FF8A000110000016FF85020101085B5D737472696E6701FF8600010C000037FF91040101266D61705B737472696E675D2A6C6F676963615F6A7565676F2E45737461646F4A756761646F7201FF9200010C01FF8C000039FF8B030102FF8C000103010643617274617301FF90000111556C74696D6F496E646963654C6569646F010400010654726F706173010400000023FF8F020101145B5D6C6F676963615F6A7565676F2E436172746101FF900001FF8E000041FF8D03010105436172746101FF8E00010401074964436172746101040001045469706F0104000106526567696F6E01040001094573436F6D6F64696E010200000014FF93020101065B5D626F6F6C01FF94000102000046FF97040101356D61705B6C6F676963615F6A7565676F2E4E756D526567696F6E5D2A6C6F676963615F6A7565676F2E45737461646F526567696F6E01FF9800010401FF96000027FF95030102FF9600010201084F637570616E7465010C0001094E756D54726F706173010400000010FF830501010454696D6501FF84000000FFD1FF88014B5A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9D03010113416363696F6E52656369626972526567696F6E01FF9E00010501084944416363696F6E0104000106526567696F6E010400010F54726F70617352657374616E74657301040001145465727269746F72696F7352657374616E74657301040001074A756761646F72010C000000FE126FFF9E0F0326015201086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102020126015001086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102040126014E01086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102060126014C01086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102080126014A01086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11020A0126014801086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11020C0124014601086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11020E0124014401086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102100124014201086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102120124014001086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102140124013E01086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102160124013C01086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102180122013A01086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11021A0122013801086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11021C0122013601086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11021E0122013401086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102200122013201086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102220122013001086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102240120012E01086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102260120012C01086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E1102280120012A01086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11022A0120012801086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11022C0120012601086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11022E0120012401086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110230011E012201086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110232011E012001086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110234011E011E01086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110236011E011C01086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110238011E011A01086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11023A011E011801086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11023C011C011601086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11023E011C011401086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110240011C011201086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110242011C011001086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110244011C010E01086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110246011C010C01086A756761646F7236005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110248011A010A01086A756761646F7231005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11024A011A010801086A756761646F7232005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11024C011A010601086A756761646F7233005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E11024E011A010401086A756761646F7234005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E110250011A010201086A756761646F7235005A6769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E52656369626972526567696F6EFF9E0F0252011A02086A756761646F723600576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FF9F03010110416363696F6E43616D62696F4661736501FFA000010301084944416363696F6E01040001044661736501040001074A756761646F72010C000000FFEDFFA00D010202086A756761646F723200586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA103010111416363696F6E496E6963696F5475726E6F01FFA200010501084944416363696F6E01040001074A756761646F72010C00010F54726F7061734F6274656E69646173010400011652617A6F6E4E756D65726F5465727269746F72696F73010400011852617A6F6E436F6E74696E656E7465734F63757061646F730104000000FE04DBFFA20D010401086A756761646F723200576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723300586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723300576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723400586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723400576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723500586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723500576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723600586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723600576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723100586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723100556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA30301010E416363696F6E5265666F727A617201FFA400010401084944416363696F6E01040001074A756761646F72010C0001135465727269746F72696F5265666F727A61646F010400010E54726F7061735265667565727A6F0104000000FE0C7CFFA40F010801086A756761646F7231021A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723200586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723200556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA411010801086A756761646F72320102011A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723300586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723300556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA411010801086A756761646F72330104011A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723400586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723400556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA411010801086A756761646F72340106011A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723500586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723500556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA411010801086A756761646F72350108011A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00D010202086A756761646F723600586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA20D010401086A756761646F723600556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA411010801086A756761646F7236010A011A00576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00F0102010201086A756761646F723100586769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E496E6963696F5475726E6FFFA211010401086A756761646F72310106010E00556769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E5265666F727A6172FFA40F010801086A756761646F7231020600576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00F0102010401086A756761646F723100576769746875622E636F6D2F554E495A41522D33303232362D323032322D30312F70726F796563746F5F736F6674776172655F6261636B656E642F6C6F676963615F6A7565676F2E416363696F6E43616D62696F46617365FFA00F0102010601086A756761646F7231000106086A756761646F7231086A756761646F7232086A756761646F7233086A756761646F7234086A756761646F7235086A756761646F72360106086A756761646F72310106011A021A00013201020132000116021600012A0102012A00011E021E0001240102012400010100086A756761646F7232020100086A756761646F7233020100086A756761646F7234020100086A756761646F7235020100086A756761646F723602010002060101010101010106012A1A01086A756761646F72320102004201086A756761646F72340102004A01086A756761646F72320102000401086A756761646F7233011C001201086A756761646F72340102002801086A756761646F72330102002A01086A756761646F72340102003401086A756761646F72330102003C01086A756761646F72310102005001086A756761646F72350102000E01086A756761646F72320102002001086A756761646F72350102003E01086A756761646F72320102000001086A756761646F72310122003001086A756761646F72310102002601086A756761646F72320102002E01086A756761646F72360102004401086A756761646F72350102000C01086A756761646F72310102001601086A756761646F72360102001001086A756761646F72330102001801086A756761646F72310102001C01086A756761646F72330102001E01086A756761646F72340102002C01086A756761646F72350102003601086A756761646F72340102000801086A756761646F7235011C000A01086A756761646F7236011C004001086A756761646F72330102003801086A756761646F72350102003A01086A756761646F72360102004C01086A756761646F72330102004E01086A756761646F72340102005201086A756761646F72360102003201086A756761646F72320102004601086A756761646F72360102001401086A756761646F72350102002201086A756761646F72360102002401086A756761646F72310102004801086A756761646F72310102000201086A756761646F7232011C000601086A756761646F7234011C000126012E0102012E00010402040001420102014200013C0102013C000140010201400001100210000114021400012C0102012C000136010201360001180218000106020600010A020A000120022000014E0104014E0001260102012600014A0104014A00013E0102013E0001300102013000014601020146000148010401480000010202020001540301000108020800011C021C00010C020C00014C0104014C00011202120001280102012800013A0102013A00013401020134000150010401500001380102013800015201040152000122022200010E020E000144010201440001560301000C06086A756761646F7231086A756761646F7232086A756761646F7233086A756761646F7234086A756761646F7235086A756761646F7236010F010000000EDA08BC0A06B3C203007800','hex'),decode('0DFF9B020102FF9C0001FF9A000013FF99030101074D656E73616A6501FF9A00000004FF9C0000','hex'),true,NULL,true,6);



INSERT INTO backend."Participa" ("ID_partida","nombreUsuario") VALUES
                                                                   (1,'jugador1'),
                                                                   (1,'jugador2'),
                                                                   (1,'jugador3'),
                                                                   (1,'jugador4'),
                                                                   (1,'jugador5'),
                                                                   (1,'jugador6');
