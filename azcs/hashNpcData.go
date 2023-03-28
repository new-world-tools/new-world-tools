package azcs

var hashNpcData = map[uint32]string{
    0x0031c4d4: "0210a_barkeep", // 3261652
    0x00965409: "13a_agnes_00", // 9851913
    0x014f8583: "12_hamzakillic", // 21988739
    0x0169ad07: "0232a_naturalist_1", // 23702791
    0x01e1aabb: "1523_fisher", // 31566523
    0x02097f56: "9514_fergus_s01_01", // 34176854
    0x033c63fe: "0907_commander", // 54289406
    0x033e121b: "1129_jaansen_note", // 54399515
    0x0345e70c: "0251a_headstone", // 54912780
    0x03639003: "9905e_msq_rima", // 56856579
    0x03e0dcff: "1037_rowanpennbrook", // 65068287
    0x0418d5f1: "12_mb_grace_04", // 68736497
    0x04893b73: "1507_commander", // 76102515
    0x04f6ea0a: "12_serpaynebennett_02", // 83290634
    0x051150aa: "1611_stranger", // 85020842
    0x05b984a8: "1128_lerato2", // 96044200
    0x06825935: "06_saatvikagrawal02", // 109205813
    0x069e6832: "0438a_jacobevans", // 111044658
    0x07a60cc0: "0437a_felixlivius", // 128322752
    0x080a0ac1: "1123_fisher", // 134875841
    0x082ed9ef: "99b_yonas_p7", // 137288175
    0x085d3e7c: "12_jeongnabi", // 140328572
    0x08aa16c5: "9508_quirinusflaccus_s01_01", // 145364677
    0x08ddf913: "1129_jaansen", // 148764947
    0x090b61f8: "12_shazaancortis", // 151740920
    0x09129546: "1686_qabalah", // 152212806
    0x09f5ffa0: "0209b_adjudicator", // 167116704
    0x0a06cf88: "0325_ranger", // 168218504
    0x0a1dd9f0: "06_margotemerson01", // 169728496
    0x0a1e7c0c: "0728_kesava", // 169769996
    0x0a2744df: "0439a_anwirhughes", // 170345695
    0x0a409469: "0250a_baldr_1", // 172004457
    0x0b156b42: "0449a_heronnote", // 185953090
    0x0b821a66: "1666_razin", // 193075814
    0x0ce166e1: "0805_artificer", // 216098529
    0x0cffdafc: "1323_fisher", // 218094332
    0x0d290d19: "0608_alchemist", // 220794137
    0x0d69bbbc: "0443a_oushen", // 225033148
    0x0e7de044: "9905c_msq_rima", // 243130436
    0x0eb9ab0b: "1417_warden", // 247048971
    0x0eecf494: "0203a_survivalist_1", // 250410132
    0x0f2820cb: "1616_niobe", // 254288075
    0x0f431df6: "99b_yonas_p3", // 256056822
    0x0f581440: "06_michaellewis03", // 257430592
    0x0f67b84e: "9808_cleo_01", // 258455630
    0x0fafee61: "dungeon_cutlasskeys00_scalleywag_00", // 263188065
    0x0fc620a4: "9201_cleric", // 264642724
    0x0fd3ef1a: "13a_pyro_01", // 265547546
    0x0fe770d4: "9805_tibbs_02", // 266825940
    0x1006be0a: "0308_alchemist", // 268877322
    0x10313508: "06_marceldupont01", // 271660296
    0x103c2f95: "0249a_aethelwynn_1", // 272379797
    0x106a4225: "99a_yonas_00_b", // 275399205
    0x10911552: "06_fancybrushset", // 277943634
    0x117d876c: "01_gracemb_00", // 293439340
    0x11a91764: "9803_morellus", // 296294244
    0x11e061b1: "1624_imhotep_v4p1", // 299917745
    0x1292ed52: "9907_msq_aevillager1", // 311618898
    0x12be4ed4: "9505_skye_s01_01", // 314461908
    0x12f52428: "9502_zander_s01_06", // 318055464
    0x13396de5: "0440a_atticusdelphi", // 322530789
    0x138ba64a: "9503_dog_s01_05", // 327919178
    0x14af7734: "1624_imhotep_v3p1", // 347043636
    0x14e66253: "9503_dog_s01_01", // 350642771
    0x156d1d03: "1624_imhotep_v2p1", // 359472387
    0x158530e1: "13a_animalhunter_03", // 361050337
    0x1598e031: "9502_zander_s01_02", // 362340401
    0x15c11160: "01_reesemb_02", // 364974432
    0x15d38acd: "9505_skye_s01_05", // 366185165
    0x167b8d7f: "1105_artificer", // 377195903
    0x168da5a8: "1624_imhotep_v4p5", // 378381736
    0x17015097: "9820_smiggs", // 385962135
    0x172ba35a: "1624_imhotep_v1p1", // 388735834
    0x1751eb8c: "0249a_aethelwynn_5", // 391244684
    0x17c3328d: "99c_bellaruby", // 398668429
    0x17c581d7: "1636_kollauthis", // 398819799
    0x18bb8656: "0247a_anora", // 414942806
    0x18f928d5: "9906_msq_shamaness", // 418982101
    0x19a68860: "1007_commander", // 430344288
    0x1a81973e: "06_royahakimi03", // 444700478
    0x1a88f34b: "1681_innkeep", // 445182795
    0x1ab89680: "13a_guard1", // 448304768
    0x1b77bca7: "0402a_magistrate", // 460831911
    0x1bac2f3e: "9918_msq_researcher4", // 464269118
    0x1bbbe7b4: "1626_hermogenes", // 465299380
    0x1c65c6e6: "9505_skye_s01_09", // 476432102
    0x1cd88fab: "99a_grace_00", // 483954603
    0x1d1c4ac4: "1122_innkeep", // 488393412
    0x1d6f70b9: "0809_adjudicator", // 493842617
    0x1dacb0a4: "12_knightsdecree", // 497856676
    0x1dc0d684: "12_oliviabarovier", // 499177092
    0x1e79b9ad: "13a_antiquarianspouse_00", // 511293869
    0x1facb489: "13a_primrose_02", // 531412105
    0x1fdb05ea: "0203a_survivalist_1a", // 534447594
    0x1fed8a00: "99c_adiana", // 535661056
    0x210f7ce1: "13a_monsterhunter_01", // 554663137
    0x2120f971: "1524_boatswain01", // 555809137
    0x21443b17: "99b_pennyfather_00", // 558119703
    0x2203c7b0: "99c_morgaine_a", // 570673072
    0x221c3b6f: "06_francisturner02", // 572275567
    0x2226be97: "9806_stroudahia_01", // 572964503
    0x22df6818: "99b_sun_p3", // 585066520
    0x22f512ea: "0329_adjudicator", // 586486506
    0x2304b151: "1630_tristan", // 587510097
    0x2305cf8e: "1682_innkeep", // 587583374
    0x236035b2: "12_oldegeiser", // 593507762
    0x2400cfc8: "0717_warden", // 604032968
    0x2434700b: "0248a_rose_1", // 607416331
    0x244cb6cb: "0607_commander", // 609007307
    0x2507f6e1: "9923_msq_wardenconstance", // 621278945
    0x256e03a9: "99c_morgaine_e", // 627966889
    0x2685e11a: "99a_adiana_00", // 646308122
    0x27038157: "06_fudehinode", // 654541143
    0x2764e3b9: "1525_angryearth", // 660923321
    0x2873c573: "0230a_ranger_1", // 678675827
    0x288e06e3: "1643_gaius", // 680396515
    0x2a59d82c: "0908_alchemist", // 710531116
    0x2a72dd34: "npc_heretic", // 712170804
    0x2a74a522: "9810_alejandro", // 712287522
    0x2a7f6a7d: "06_verildisjannsen", // 712993405
    0x2b2145fa: "06_kathrijnjanssen02", // 723600890
    0x2b4d0127: "1642_senhyris", // 726466855
    0x2b58ce78: "0623_fisher", // 727240312
    0x2be52191: "dungeon_brokenlever", // 736436625
    0x2be85a71: "13a_architect_03", // 736647793
    0x2bed9f74: "01_graceww_01", // 736993140
    0x2c9830b9: "1415_watcher", // 748171449
    0x2d877366: "0205a_artificer", // 763851622
    0x2deb6fb2: "0317_warden", // 770404274
    0x2dec80a1: "1508_alchemist", // 770474145
    0x2e00ac5f: "7505_artificer", // 771796063
    0x2e0e8665: "12_percillawallen", // 772703845
    0x2e6f7909: "1001_watcher", // 779057417
    0x2edf1ad7: "0715_watcher", // 786373335
    0x2f1f8274: "06_graceww_04", // 790594164
    0x2f510978: "01_reeseww_03", // 793839992
    0x2fad1e45: "0423_fisher", // 799874629
    0x302af3b3: "1625_ghazi_p1", // 808121267
    0x3097b396: "9504_shailyn_s01_01", // 815248278
    0x30c333b2: "1008_alchemist", // 818099122
    0x31a09caa: "9825_pop", // 832609450
    0x31abb994: "9501_grace_s01_06", // 833337748
    0x31accfa9: "9811_gregario", // 833408937
    0x31b65a31: "0811_stranger", // 834034225
    0x3202ab3f: "1645_xiaoling", // 839035711
    0x322563d1: "06_towerjournal", // 841311185
    0x32feffa5: "0804_overseer", // 855572389
    0x33024cbe: "0831_tangzhi_start", // 855788734
    0x34275f47: "1622_innkeep", // 874995527
    0x3484d19a: "99b_galahad_p4", // 881119642
    0x34fe2f2d: "reekwater_dungeon_scylla_01", // 889073453
    0x357c9b36: "12_regentjinjae", // 897358646
    0x358c60ca: "0252a_guard", // 898392266
    0x3591fecb: "0823_fisher", // 898760395
    0x35b3e7f5: "12_lochnir", // 900982773
    0x36607529: "0522_innkeep", // 912291113
    0x366fd26d: "99b_morgengrave_c", // 913298029
    0x36b7304e: "06_leovixsilva01", // 917975118
    0x36c67d8d: "9501_grace_s01_02", // 918977933
    0x36d46272: "7400_lorenotetest", // 919888498
    0x374737aa: "1625_ghazi_p5", // 927414186
    0x37bbebb8: "1679_tomash_er", // 935062456
    0x37c1a117: "9204_executioner", // 935436567
    0x3888e796: "1606_ambrosius", // 948496278
    0x38fa9a81: "0439a_anwirhughes_02", // 955947649
    0x396305d8: "0307_commander", // 962790872
    0x398d6855: "0202a_magistrate", // 965568597
    0x39a40e94: "0714_watcher", // 967052948
    0x39a864f8: "12_mb_grace", // 967337208
    0x3a260f3b: "1634_phimenis", // 975572795
    0x3a2ff718: "9511_dicunhobday_s01_02", // 976221976
    0x3a53a841: "0245a_liang", // 978561089
    0x3af15489: "9932_vk_romacnote", // 988894345
    0x3afe69b2: "1649_miresis", // 989751730
    0x3b163db4: "0326_prospector", // 991313332
    0x3bb13620: "1655_sokmenis", // 1001469472
    0x3be324fa: "1414_watcher", // 1004741882
    0x3be9c5e6: "06_nessaharrower03", // 1005176294
    0x3c4afc55: "12_percillawallen_03", // 1011547221
    0x3c7ae7ee: "0922_innkeep", // 1014687726
    0x3c8401ff: "06_nessaharrower07", // 1015284223
    0x3cd75869: "9919_msq_wardenscout", // 1020745833
    0x3d329db1: "99b_galahad_p8", // 1026727345
    0x3d3399e5: "1605_akila_corsica", // 1026791909
    0x3df291d3: "1671_commander2", // 1039307219
    0x3e361dc6: "99c_galahad_d", // 1043733958
    0x3e45366f: "1510_barkeep", // 1044723311
    0x3f02b88c: "1660_paubastis", // 1057142924
    0x3f6fe7f7: "9536_noraroggeveen", // 1064298487
    0x3fc81a2f: "13a_antiquarian_00", // 1070078511
    0x400401a5: "1033_treasurehunter", // 1074004389
    0x4040073c: "1625_ghazi_p4", // 1077937980
    0x404c8ca4: "0448a_syndicatealchemist", // 1078758564
    0x40ae45f7: "0831_tangzhi_final", // 1085162999
    0x412ef70a: "12_hugomolina", // 1093596938
    0x4168e2fb: "99b_morgengrave_b", // 1097392891
    0x42e6c1ec: "0320_ranger", // 1122419180
    0x434a8948: "9917_msq_researcher3", // 1128958280
    0x4383e10c: "99b_galahad_p5", // 1132716300
    0x439bf714: "0802_magistrate", // 1134294804
    0x43efb664: "0249a_aethelwynn_2a", // 1139783268
    0x4494dba2: "reekwater_dungeon_scylla_04", // 1150606242
    0x44ee2515: "99b_galahad_p1", // 1156457749
    0x45111970: "9902_msq_adiana", // 1158748528
    0x4570f305: "1322_innkeep", // 1165030149
    0x46038a32: "06_tsuifen", // 1174637106
    0x4651df06: "1625a_ghazi_p8", // 1179770630
    0x4667323c: "1673_adjudicator2", // 1181168188
    0x46ac8902: "9501_grace_s01_07", // 1185712386
    0x46ddc4c1: "06_leovixsilva04", // 1188938945
    0x47347d83: "9904d_msq_theophrastus", // 1194622339
    0x477622bb: "1419_ranger", // 1198924475
    0x477d2a23: "1599_nekumanesh", // 1199385123
    0x48cf2ab9: "13a_antiquarian_01", // 1221536441
    0x48f95d70: "1412_watcher", // 1224301936
    0x49312d50: "99c_galahad_e", // 1227959632
    0x4948c106: "9904b_msq_theophrastus", // 1229504774
    0x49c78ecd: "06_francisturner", // 1237814989
    0x49e7b324: "0709_adjudicator", // 1239921444
    0x4a87b47c: "7404_overseer", // 1250407548
    0x4abe771e: "0712_watcher", // 1253996318
    0x4b0d6196: "0720_ranger", // 1259168150
    0x4b4dccc3: "12_percillawallen_02", // 1263389891
    0x4b7362ac: "0425_herbalistmb", // 1265853100
    0x4b833169: "06_nessaharrower06", // 1266889065
    0x4b98bd41: "9929_msq_recliningthorpe", // 1268301121
    0x4b9b88a2: "0207a_commander", // 1268484258
    0x4bcd069f: "0505_artificer", // 1271727775
    0x4c1047df: "9910_msq_aescout", // 1276135391
    0x4cac9a7e: "9506_granny_s01_02", // 1286380158
    0x4ceef570: "06_nessaharrower02", // 1290728816
    0x4d198f3d: "9908_msq_aevillager2", // 1293520701
    0x4d2c207f: "1678_emptychest", // 1294737535
    0x4d4bd040: "12_shaegarcia", // 1296814144
    0x4d9a5ca7: "1502_magistrate", // 1301961895
    0x4e25b435: "06_gideonstormsbane02", // 1311093813
    0x4e5ce949: "99c_galahad_a", // 1314711881
    0x4e97ca8a: "1628_bekis", // 1318570634
    0x4f099324: "1509_adjudicator", // 1326027556
    0x4f129a2d: "1010_barkeep", // 1326619181
    0x4f65d04c: "06_yelenapajitnova01", // 1332072524
    0x4fcb40d2: "06_medea01", // 1338720466
    0x500c0a11: "1672_alchemist2", // 1342966289
    0x51597bbf: "9814_beyza", // 1364818879
    0x522759d8: "1653_lavinia", // 1378310616
    0x5269333f: "99c_morgaine_d", // 1382626111
    0x537c0781: "0447a_locket", // 1400637313
    0x5393d308: "1675_alchemist3", // 1402196744
    0x53bbb49c: "1038_quiggleyquorus", // 1404810396
    0x5436d3bc: "1102_magistrate", // 1412879292
    0x55218e01: "9806_stroudahia_00", // 1428262401
    0x5524605f: "0725_prospector", // 1428447327
    0x55d8588e: "99b_sun_p2", // 1440241806
    0x55d87c05: "1225_giacomotravelersb", // 1440250885
    0x56084c77: "13a_monsterhunter_00", // 1443384439
    0x5627c9e7: "1524_boatswain00", // 1445448167
    0x56fe3e31: "9824_kuro", // 1459502641
    0x5706f30b: "0451a_minecart", // 1460073227
    0x5801792f: "99a_captain", // 1476491567
    0x5818b2e2: "06_graceww_05", // 1478013666
    0x585639ee: "01_reeseww_02", // 1482045934
    0x585c724f: "1103_survivalist", // 1482453583
    0x58f53161: "9821_syndicateguard_m", // 1492463969
    0x59252a75: "1124_pierre_note", // 1495607925
    0x59463c31: "9809_covenantguard_f", // 1497775153
    0x5ad9ca6e: "9921_msq_wardenwalker", // 1524222574
    0x5b4fb265: "06_talesaserrano", // 1531949669
    0x5b6e4f94: "1124_pierreauguste2", // 1533955988
    0x5b86a0c1: "1605_akila_ennead", // 1535549633
    0x5bfb663b: "0405a_artificerpetrowski_01", // 1543202363
    0x5c26756c: "06_kathrijnjanssen03", // 1546024300
    0x5c55fbe9: "99a_anthurus", // 1549138921
    0x5cdd06b9: "06_lizihan", // 1557989049
    0x5ceaafe2: "01_graceww_00", // 1558884322
    0x5cef6ae7: "13a_architect_02", // 1559194343
    0x5dc5635d: "0713_watcher", // 1573217117
    0x5df63849: "1648_flora", // 1576417353
    0x5e74f95d: "1609_adjudicator", // 1584724317
    0x5f1b8e2e: "9926_vk_abigailroseb", // 1595641390
    0x5f27d42d: "0242a_ede", // 1596445741
    0x5f7576fb: "06_graceww_01", // 1601533691
    0x5f824933: "1413_watcher", // 1602373939
    0x5fd3f247: "9903b_msq_emile", // 1607725639
    0x606e27a2: "12_clericletter", // 1617831842
    0x608c3322: "9203_firepriest", // 1619800866
    0x618a953e: "1624_imhotep_v4p4", // 1636472126
    0x61a469f4: "12_allyamusa", // 1638164980
    0x61b40555: "1420_ranger", // 1639187797
    0x621a913d: "1650_flame", // 1645908285
    0x627082cb: "1225_giacomo_mb", // 1651540683
    0x62820077: "13a_animalhunter_02", // 1652686967
    0x6293b6e2: "9501_grace_s01_03a", // 1653847778
    0x629bcfac: "0441a_sairayangandul", // 1654378412
    0x629fd0a7: "9502_zander_s01_03", // 1654640807
    0x62c621f6: "01_reesemb_03", // 1657151990
    0x62d4ba5b: "9505_skye_s01_04", // 1658108507
    0x630ad966: "99a_yonas_02", // 1661655398
    0x63fe90c1: "1528_rogue", // 1677627585
    0x6424e602: "0319_ranger", // 1680139778
    0x648c96dc: "9503_dog_s01_04", // 1686935260
    0x64b5aa0b: "12_jeongiseul", // 1689627147
    0x64d0203d: "12_regentjinjae_02", // 1691361341
    0x65b1dcf8: "0435_announcer", // 1706155256
    0x663c7607: "1225_giacomo_fl", // 1715238407
    0x667ab7fa: "01_gracemb_01", // 1719318522
    0x66a1a556: "1601_gideon", // 1721869654
    0x66b83e51: "0436a_tamraayad", // 1723350609
    0x66d4fa4f: "1637_livia", // 1725233743
    0x674157d5: "1624_imhotep_v1p4", // 1732335573
    0x679523e7: "9928_msq_yonasmemorial", // 1737827303
    0x67bc9935: "0403a_huntsmanlee", // 1740413237
    0x683a8ea8: "0234a_mourner_1", // 1748668072
    0x697e893b: "13a_antiquarianspouse_01", // 1769900347
    0x69a609c5: "99c_grace", // 1772489157
    0x69f25962: "9507_williamdaintith_s01_01", // 1777490274
    0x69f85d77: "1605_akila_menefir", // 1777884535
    0x6a0c220a: "0404a_constable", // 1779180042
    0x6b62f670: "9505_skye_s01_08", // 1801647728
    0x6bdbccc4: "1403_survivalist", // 1809566916
    0x6bdfbf3d: "99a_grace_01", // 1809825597
    0x6c4be686: "1422_innkeep", // 1816913542
    0x6cdbab5a: "0445a_requisitionnote", // 1826335578
    0x6cdf8e9d: "1664_aban", // 1826590365
    0x6d11ef0d: "0246a_cornelius", // 1829891853
    0x6d121fda: "12_geertielother", // 1829904346
    0x6d86a7a8: "06_royahakimi02", // 1837541288
    0x6dcf4678: "0719_ranger", // 1842300536
    0x6df347d6: "1309_adjudicator", // 1844660182
    0x6e0ccce8: "0722_innkeep", // 1846332648
    0x6e743707: "0810_barkeep", // 1853110023
    0x6efc748c: "7507_swap", // 1862038668
    0x6f3da9c8: "1679_tomash_gc", // 1866312136
    0x6fb34a09: "1640_steffan", // 1874020873
    0x6fc5a594: "1526_naturalist", // 1875223956
    0x7001ba67: "0425_herbalistmb_town", // 1879161447
    0x70278b3a: "12_adjudicatorowgan", // 1881639738
    0x7038aa67: "9930_vk_romacbarboa", // 1882761831
    0x70914b3b: "99b_grace_p1", // 1888570171
    0x718569a3: "06_saatvikagrawal03", // 1904568739
    0x7189b925: "1111_stranger", // 1904851237
    0x71c727ce: "1504_overseer", // 1908877262
    0x7236f1d0: "1632_titus", // 1916203472
    0x72c11cb1: "1104_overseer", // 1925258417
    0x73f1da9c: "12_serpaynebennett_03", // 1945229980
    0x7472217e: "12_mb_grace_01", // 1953636734
    0x748b32aa: "1606_ambrosius_p2", // 1955279530
    0x74d56326: "1403_survivalista", // 1960141606
    0x768d3c9c: "1603_barnabas_coronation", // 1988967580
    0x76de31a4: "1527_wader", // 1994273188
    0x76e8e7c4: "9903e_msq_emile", // 1994975172
    0x78442d60: "99b_yonas_p2", // 2017733984
    0x785f24d6: "06_michaellewis02", // 2019501270
    0x786088d8: "9808_cleo_00", // 2019592408
    0x78a8def7: "dungeon_cutlasskeys00_scalleywag_01", // 2024333047
    0x78ba2dc9: "9801_florina_02", // 2025467337
    0x78d4df8c: "13a_pyro_00", // 2027216780
    0x7977d8ab: "0723_innkeep", // 2037897387
    0x797878a3: "9915_msq_researcher1", // 2037938339
    0x7a10af70: "0428_jones", // 2047913840
    0x7a600daf: "0209_adjudicator", // 2053115311
    0x7b30f2c5: "1423_innkeep", // 2066805445
    0x7b376383: "06_lioneldelarue02", // 2067227523
    0x7b501371: "0408a_alchemist", // 2068845425
    0x7b77e8c2: "06_zhanghaoyu", // 2071455938
    0x7bd311c3: "1615_nicabar", // 2077430211
    0x7c4886bd: "0503_survivalist", // 2085127869
    0x7c8e2daf: "1009_adjudicator", // 2089692591
    0x7ca24f03: "9505_skye_s01_10", // 2091011843
    0x7d350ba4: "1030_tamer", // 2100628388
    0x7db5aa95: "13a_niall_00", // 2109057685
    0x7e4d3884: "9815_marauderguard_m", // 2118989956
    0x7e5b236f: "06_lijsbetdevries02", // 2119902063
    0x7e91ddd2: "1661_jackal", // 2123488722
    0x7edc76bc: "0223a_fisher", // 2128377532
    0x7f1dc489: "13a_maria_00", // 2132657289
    0x7f29e979: "99b_yonas_p6", // 2133453177
    0x8046f4b1: "0803_survivalist", // 2152133809
    0x82228448: "1644_nesse_outpost", // 2183300168
    0x8293ffcd: "1629_elissa", // 2190737357
    0x829674f7: "1031_swampwalker", // 2190898423
    0x83848a60: "0244a_sebastion", // 2206501472
    0x839258c7: "0208a_alchemist", // 2207406279
    0x83b1c73a: "13a_guard2", // 2209466170
    0x86a5e533: "13a_primrose_01", // 2259019059
    0x87a76108: "0721_ranger", // 2275893512
    0x87fcba4a: "13a_perfumealchemist", // 2281486922
    0x881d2847: "npc_fatherrusso", // 2283612231
    0x8874d6d6: "01_gracemb_03", // 2289358550
    0x88ccb053: "9813_damiano", // 2295115859
    0x88e9300b: "1624_imhotep_v4p2", // 2296983563
    0x89357e2f: "0249a_aethelwynn_2", // 2301984303
    0x893864b2: "06_marceldupont02", // 2302174386
    0x8963139f: "99a_yonas_00_a", // 2304971679
    0x897d7d16: "0314_watcher", // 2306702614
    0x8a0e6d75: "1603_barnabas", // 2316201333
    0x8a31d573: "1647_aulus", // 2318521715
    0x8a82f7f0: "9503_dog_s01_06", // 2323838960
    0x8b45e842: "0438a_jacobevans_01", // 2336614466
    0x8b5229c0: "0899_liberator", // 2337417664
    0x8bb71f6e: "9505_skye_s01_02", // 2344034158
    0x8bdc2225: "1418_ranger", // 2346459685
    0x8bfc7592: "9502_zander_s01_05", // 2348578194
    0x8c26e579: "1617_publius", // 2351359353
    0x8c3e2152: "1225_giacomochason", // 2352882002
    0x8c644cb9: "1624_imhotep_v2p2", // 2355383481
    0x8c8c615b: "13a_animalhunter_00", // 2358010203
    0x8c91b18b: "9502_zander_s01_01", // 2358358411
    0x8cc840da: "01_reesemb_01", // 2361934042
    0x8cdadb77: "9505_skye_s01_06", // 2363153271
    0x8d04b84a: "99a_yonas_00", // 2365896778
    0x8d095f7c: "1308_alchemist", // 2366201724
    0x8dc60429: "9905_msq_rima", // 2378564649
    0x8def33e9: "9503_dog_s01_02", // 2381263849
    0x8e22f2e0: "1624_imhotep_v1p2", // 2384655072
    0x8e4cc172: "0321_ranger", // 2387394930
    0x8e58ba36: "0249a_aethelwynn_6", // 2388179510
    0x8e9c45ed: "1110_barkeep", // 2392606189
    0x8eee41a4: "0824_innkeep", // 2397979044
    0x8f3eebb8: "9902_bsmsq_adiana", // 2403265464
    0x8f7d961b: "0699_barkimedes", // 2407372315
    0x9026ec6f: "1608_alchemist", // 2418469999
    0x90be9e6c: "1425_prospector", // 2428411500
    0x91278855: "99b_yonas_p4", // 2435287125
    0x917e8a9a: "1645_xiaoling_cavalier_rest", // 2440989338
    0x91941a3e: "1683_zahur", // 2442402366
    0x9314884a: "06_margotemerson02", // 2467596362
    0x9349c5d3: "0250a_baldr_2", // 2471085523
    0x93f7b63a: "9202_hunter", // 2482484794
    0x940f8844: "0446a_ancientscroll", // 2484045892
    0x94607cdd: "12_innkeepbiton", // 2489351389
    0x94da6506: "0231a_painter", // 2497340678
    0x96837c9c: "0425_herbalistww", // 2525199516
    0x968c8126: "9206_guardian", // 2525790502
    0x96b44ce5: "9801_florina_00", // 2528398565
    0x96c19af8: "9819_hedge", // 2529270520
    0x96dabea0: "13a_pyro_02", // 2530918048
    0x96ee216e: "9805_tibbs_01", // 2532188526
    0x97e5a52e: "0203a_survivalist_2", // 2548409646
    0x985aa998: "9512_marty_s01_01", // 2556078488
    0x9860fcbd: "0232a_naturalist_2", // 2556492989
    0x9891c47e: "99b_yonas_p8", // 2559689854
    0x9986da05: "0507_commander", // 2575751685
    0x998afc6b: "0001_captainthorpe", // 2576022635
    0x9a26ff83: "1613_alaina", // 2586247043
    0x9a7c4052: "12_mb_grace_03", // 2591834194
    0x9b002eec: "9514_fergus_s01_02", // 2600480492
    0x9b933ede: "9812_peta", // 2610118366
    0x9c857974: "12_hugomolina_01", // 2625993076
    0x9de4d0fb: "9924_msq_wardeneliza", // 2649018619
    0x9dffbbb0: "12_serpaynebennett_01", // 2650782640
    0x9e066955: "0315_watcher", // 2651220309
    0x9e40648a: "1654_margarete", // 2655020170
    0x9e5d2d87: "9903c_msq_emile", // 2656906631
    0x9e63222a: "1129_jaansen2", // 2657296938
    0x9f560d05: "dungeon_cutlasskeys00_mariner_01", // 2673216773
    0x9f8b088f: "06_saatvikagrawal01", // 2676689039
    0xa025ca68: "9207_rogue", // 2686831208
    0xa0528865: "99c_galahad_c", // 2689763429
    0xa06b0af8: "9902b_msq_adiana", // 2691369720
    0xa114a405: "9501_grace_s01_09", // 2702484485
    0xa153044a: "06_ahurafarid", // 2706572362
    0xa1560bf2: "12_gaetanfortier", // 2706770930
    0xa16546e6: "0718_ranger", // 2707769062
    0xa16bb160: "06_yelenapajitnova03", // 2708189536
    0xa1c521fe: "06_medea03", // 2714051070
    0xa1f3cb3b: "0439a_anwirhughes_01", // 2717109051
    0xa23dc7c9: "9802_lethold", // 2721957833
    0xa279d41c: "9926_vk_abigailroseep2", // 2725893148
    0xa326a6a2: "9511_dicunhobday_s01_01", // 2737219234
    0xa335658f: "1663_esoeris", // 2738185615
    0xa423d23b: "12_fionamurphy", // 2753811003
    0xa46ce4ae: "1307_commander", // 2758599854
    0xa494ab51: "0909_adjudicator", // 2761206609
    0xa4e05fc4: "9904_msq_theophrastus", // 2766168004
    0xa58d5045: "06_nessaharrower04", // 2777501765
    0xa5e07a00: "0510_barkeep", // 2782951936
    0xa5f0df92: "0832_alvaro_start", // 2784026514
    0xa61a77e7: "9904c_msq_theophrastus", // 2786752487
    0xa62b539a: "1679_tomash_sm", // 2787857306
    0xa6c14b95: "13a_antiquarian_03", // 2797685653
    0xa7a7506e: "1610_barkeep", // 2812760174
    0xa81b0399: "0727_bookofisabella", // 2820342681
    0xa866cb62: "9904e_msq_theophrastus", // 2825309026
    0xa8705289: "12_dayomusa", // 2825933449
    0xa88ee69c: "0318_ranger", // 2827937436
    0xa8a2e82e: "9501_grace_s01_05", // 2829248558
    0xa8b1d10c: "0450a_townguard", // 2830225676
    0xa923a209: "1625_ghazi_p2", // 2837684745
    0xa9cbd6d2: "0430_pvparena", // 2848708306
    0xaa2c1352: "1602_crassus_scout", // 2855015250
    0xaa55e425: "9903_msq_emile", // 2857755685
    0xaab428f9: "1604_charmion", // 2863933689
    0xaae04439: "99b_galahad_p3", // 2866824249
    0xac3b1c6e: "06_nessaharrower08", // 2889555054
    0xac51fe4b: "0504_overseer", // 2891054667
    0xad1e05cb: "1421_ranger", // 2904425931
    0xad688f04: "0411a_williamheron", // 2909310724
    0xad8d8020: "99b_galahad_p7", // 2911731744
    0xada2da8d: "0243a_niko", // 2913131149
    0xadc4ccda: "1225_giacomo_ef", // 2915355866
    0xadc53946: "1522_innkeep", // 2915383622
    0xadf77e97: "reekwater_dungeon_scylla_02", // 2918678167
    0xae4e6610: "1625_ghazi_p6", // 2924373520
    0xae50264c: "9517_otmarwinkler_s01_01", // 2924488268
    0xaf195bdf: "0511_stranger", // 2937674719
    0xaf4f5856: "0440a_atticusdelphi_02", // 2941212758
    0xafbe61f4: "06_leovixsilva02", // 2948489716
    0xafcf2c37: "9501_grace_s01_01", // 2949590071
    0xb0e361d7: "0508_alchemist", // 2967691735
    0xb17a94c9: "0230a_ranger_2", // 2977600713
    0xb17b17d7: "06_graceww_03", // 2977634263
    0xb2281440: "06_kathrijnjanssen01", // 2988971072
    0xb2e10bcb: "13a_architect_00", // 3001093067
    0xb2e4cece: "01_graceww_02", // 3001339598
    0xb359c314: "0407a_commander", // 3009004308
    0xb36a61df: "9509_adolfodeacutis_s01_01", // 3010093535
    0xb3eb729d: "9916_msq_researcher2", // 3018551965
    0xb43f936e: "0209a_adjudicator", // 3024065390
    0xb4b984bd: "1646_froderico", // 3032057021
    0xb5aef67e: "1680_buhawi", // 3048142462
    0xb5c36cc7: "1676_adjudicator3", // 3049483463
    0xb66ac8b5: "9807_aventus_01", // 3060451509
    0xb6b9cbd3: "1128_lerato", // 3065629651
    0xb7663804: "9903d_msq_emile", // 3076929540
    0xb8062d5b: "13a_monsterhunter_02", // 3087412571
    0xb833b479: "1614_guard", // 3090396281
    0xb9153275: "9901_msq_yonas", // 3105174133
    0xb92c5bc2: "1641_eudoxia", // 3106692034
    0xb94357bd: "1607_commander", // 3108198333
    0xb94bd327: "0233a_hunter", // 3108754215
    0xba0eb57f: "0238_alchemist", // 3121526143
    0xbb0a960a: "99c_morgaine_b", // 3138033162
    0xbb1452e9: "1127_harunobu", // 3138671337
    0xbcbbfdbb: "99b_sun_p4", // 3166436795
    0xbdb874e8: "9909_msq_aevillager3", // 3182982376
    0xbf809cfb: "1225_giacomo_ww", // 3212877051
    0xbf8bf9ab: "1124_pierreauguste", // 3213621675
    0xc06f757b: "0523_fisher", // 3228530043
    0xc111e358: "06_graceww_06", // 3239175000
    0xc15385e5: "1602_crassus_pavilion", // 3243476453
    0xc15f6854: "01_reeseww_01", // 3244255316
    0xc16df823: "9807_aventus_00", // 3245209635
    0xc1856d31: "1602_crassus_desert", // 3246746929
    0xc1d32de0: "9501_grace_s01_10", // 3251842528
    0xc2095e4b: "9205_blademaster", // 3255393867
    0xc213ab95: "1416_warden", // 3256069013
    0xc21682d7: "9510_eynonvoyle_s01_01", // 3256255191
    0xc24f8ab1: "1602_crassus_settlement", // 3259992753
    0xc2c407ef: "9513_amis_s01_01", // 3267626991
    0xc2f23781: "0405a_artificerpetrowski_02", // 3270653825
    0xc39bd647: "0309_adjudicator", // 3281770055
    0xc3b8f083: "1127_harunobu2", // 3283677315
    0xc43051e5: "1107_commander", // 3291501029
    0xc4aedf8f: "12_alchemisttuit", // 3299794831
    0xc575f647: "1109_adjudicator", // 3312842311
    0xc5b4be54: "12_chanduiyer", // 3316956756
    0xc5e3fe58: "01_graceww_03", // 3320053336
    0xc5e63b5d: "13a_architect_01", // 3320200029
    0xc5f3fbd4: "9817_mersha", // 3321101268
    0xc612df94: "9926_vk_abigailrosea", // 3323125652
    0xc67c2741: "06_graceww_02", // 3330025281
    0xc67da45f: "0230a_ranger_3", // 3330122847
    0xc69de88e: "npc_9902_msq_yonas", // 3332237454
    0xc70f3ae3: "12_tomiradudek", // 3339664099
    0xc711eeb1: "1665_heron", // 3339841201
    0xc7617aa2: "1638_suke", // 3345054370
    0xc7912c12: "1656_perrin", // 3348179986
    0xc7b40c25: "0324_fisher", // 3350465573
    0xc7f6ade1: "99c_fritzcromer", // 3354832353
    0xc875288d: "1662_rima", // 3363121293
    0xc8af1a7b: "9208_shaman", // 3366918779
    0xc8b9cfb2: "1036_wallywiddershins", // 3367620530
    0xc93932ca: "1002_magistrate", // 3375968970
    0xc94c9f4b: "99c_ruxandracelrau", // 3377241931
    0xc9583abf: "12_youngegeiser", // 3378002623
    0xc96fbc6b: "1657_tamouthis", // 3379543147
    0xc9aeab29: "0323_innkeep", // 3383667497
    0xca3020df: "1685_acario", // 3392151775
    0xcb67ee44: "0222a_innkeep", // 3412586052
    0xcbed54fa: "1005_artificer", // 3421328634
    0xcc0da69c: "99c_morgaine_c", // 3423446684
    0xccd10934: "99b_sun_p1", // 3436251444
    0xccd12dbf: "1225_giacomotravelersa", // 3436260799
    0xcd0c8956: "12_serpaynebennett", // 3440150870
    0xcd90678f: "0208_alchemist", // 3448792975
    0xce5fac5f: "0724_fisher", // 3462376543
    0xce94e5b9: "9809_covenantguard_m", // 3465864633
    0xcf011dcd: "13a_monsterhunter_03", // 3472956877
    0xcf27e8e9: "9821_syndicateguard_f", // 3475499241
    0xcfd93ac5: "0405a_artificerpetrowski", // 3487120069
    0xd03560c7: "9818_doc", // 3493159111
    0xd0bfd49c: "0708_alchemist", // 3502232732
    0xd0d2f69a: "0524_stranger", // 3503486618
    0xd122f813: "7504_overseer", // 3508729875
    0xd1c67b03: "13a_antiquarian_02", // 3519445763
    0xd220172c: "1503_survivalist", // 3525318444
    0xd2449d79: "12_percillawallen_01", // 3527712121
    0xd256d655: "1674_commander3", // 3528906325
    0xd28a60d3: "06_nessaharrower05", // 3532284115
    0xd29472d8: "1652_lysandra", // 3532944088
    0xd5a5cbc4: "9506_granny_s01_01", // 3584412612
    0xd5ac7c1c: "12_donquixolay", // 3584850972
    0xd5e7a4ca: "06_nessaharrower01", // 3588728010
    0xd6139493: "9501_grace_s01_08", // 3591607443
    0xd65cd6db: "0832_alvaro_final", // 3596408539
    0xd66c81f6: "06_yelenapajitnova02", // 3597435382
    0xd6c21168: "06_medea02", // 3603042664
    0xd6c2e7e9: "1505_artificer", // 3603097577
    0xd6f0fc2c: "1310_barkeep", // 3606117420
    0xd72ce58f: "06_gideonstormsbane01", // 3610043791
    0xd755b8f3: "99c_galahad_b", // 3612719347
    0xd76ef9d2: "1126_leilani", // 3614374354
    0xd7cafe71: "06_elricchapman", // 3620404849
    0xd861b341: "99b_morgengrave_a", // 3630281537
    0xd8b95162: "06_leovixsilva03", // 3636023650
    0xd91fe2f6: "1407_commander", // 3642745590
    0xd9495686: "1625_ghazi_p7", // 3645462150
    0xda5395f5: "0923_fisher", // 3662910965
    0xda8ab0b6: "99b_galahad_p6", // 3666522294
    0xdac66743: "1658_mihri", // 3670435651
    0xdaf04e01: "reekwater_dungeon_scylla_03", // 3673181697
    0xdc64f93a: "1644_nesse", // 3697604922
    0xdc929504: "1022_innkeep", // 3700593924
    0xdca8f810: "1618_caeso", // 3702061072
    0xdcbdedda: "1034_navigator", // 3703434714
    0xddc86cd9: "1627_mevia", // 3720899801
    0xdde774af: "99b_galahad_p2", // 3722933423
    0xde0d7d8d: "1639_leo", // 3725426061
    0xde122b36: "9804_feldan", // 3725732662
    0xde24929f: "1625_ghazi_p3", // 3726938783
    0xdeaaba7b: "0807_commander", // 3735730811
    0xded5bf6a: "0322_innkeep", // 3738550122
    0xdfa5d8b8: "9501_grace_s01_04", // 3752188088
    0xe02e6276: "1605_akila_coronation", // 3761136246
    0xe0796388: "1687_clarisse", // 3766051720
    0xe0b4eacf: "0425_herbalistww_town", // 3769952975
    0xe1416f2c: "0316_warden", // 3779161900
    0xe14d7cda: "99b_yonas_p1", // 3779951834
    0xe156756c: "06_michaellewis01", // 3780539756
    0xe16102b5: "0509_adjudicator", // 3781231285
    0xe1a7a9a7: "1003_survivalist", // 3785861543
    0xe1b37c73: "9801_florina_01", // 3786636403
    0xe1bf8b7a: "9905b_msq_rima", // 3787426682
    0xe1e911f8: "9805_tibbs_00", // 3790148088
    0xe2176acc: "1125_mckeys", // 3793185484
    0xe23e3239: "06_lioneldelarue01", // 3795726905
    0xe2682e25: "0444a_lostalchemy", // 3798478373
    0xe3098fbc: "1635_antonia", // 3809054652
    0xe3346de0: "06_priyathakur", // 3811864032
    0xe33db1c2: "1223_fisher", // 3812471234
    0xe41226a2: "1633_murdoch", // 3826394786
    0xe4e6c89c: "1424_fisher", // 3840329884
    0xe4f5dc5d: "0207_commander", // 3841317981
    0xe56e95f3: "9826_astrid", // 3849229811
    0xe5d009d7: "9925_msq_wardensamir", // 3855616471
    0xe620b8c3: "99b_yonas_p5", // 3860904131
    0xe661b045: "12_lochnir_01", // 3865161797
    0xe684da4d: "12_ashedubois", // 3867466317
    0xe708fa7a: "1631_tasechmi", // 3876125306
    0xe75272d5: "06_lijsbetdevries01", // 3880940245
    0xe77fcbf0: "9922_msq_wardenwayne", // 3883912176
    0xe7aefa8a: "0409a_adjudicator", // 3887004298
    0xe7c861ff: "1023_fisher", // 3888669183
    0xe7ef7391: "1129_sculpture", // 3891229585
    0xe815c915: "1612_andrea", // 3893741845
    0xe83165de: "0726_mararosa", // 3895551454
    0xe8513d93: "dungeon_cutlasskeys00_mariner_00", // 3897638291
    0xe8893567: "7506_spawn", // 3901306215
    0xe8aacf56: "0716_warden", // 3903508310
    0xe931fbdd: "9930_vk_romacbarbob", // 3912367069
    0xe96450de: "1004_overseer", // 3915665630
    0xe9981a81: "99b_grace_p2", // 3919059585
    0xe99be308: "1605_akila_marcella", // 3919307528
    0xe99fe10c: "9815_marauderguard_f", // 3919569164
    0xea2cf54a: "1011_stranger", // 3928814922
    0xead611b8: "1623_fisher", // 3939897784
    0xeb330c3d: "9816_pavel", // 3945991229
    0xebd12990: "9931_vk_caoimhe", // 3956353424
    0xec602c53: "1684_tahirah", // 3965725779
    0xeca1fb3d: "9905d_msq_rima", // 3970038589
    0xed1c10df: "0313_watcher", // 3978039519
    0xed454bd8: "9914_msq_hunter", // 3980741592
    0xed55ea37: "1108_alchemist", // 3981830711
    0xed7b70c4: "12_mb_grace_02", // 3984289988
    0xedf9e60c: "1601_gideon_coronation", // 3992577548
    0xee4243c9: "0411a_williamheron_01", // 3997320137
    0xef67cc2b: "0232a_naturalist_3", // 4016557099
    0xef96f4e8: "99b_yonas_p9", // 4019647720
    0xf01c68cc: "0609_adjudicator", // 4028393676
    0xf07a5924: "1408_alchemist", // 4034550052
    0xf0edc232: "12_eliasderit", // 4042113586
    0xf1304fc8: "12_cemalqadir", // 4046475208
    0xf133df12: "0234a_mourner_2", // 4046708498
    0xf1a2d5a5: "13a_primrose_00", // 4053980581
    0xf1d3b5b8: "0442a_concettasacco", // 4057183672
    0xf205f35f: "06_utkarshsingh", // 4060476255
    0xf39675da: "7405_artificer", // 4086724058
    0xf461874f: "0502_magistrate", // 4100032335
    0xf48ff612: "06_royahakimi01", // 4103075346
    0xf5347af9: "1651_enez", // 4113857273
    0xf59d16ff: "1027_commander", // 4120712959
    0xf6f248cc: "1409_adjudicator", // 4143073484
    0xf7cf01a9: "0808_alchemist", // 4157538729
    0xf7eeaae9: "0422_innkeep", // 4159613673
    0xf925c276: "1624_imhotep_v1p3", // 4180001398
    0xf9266054: "0253a_guard", // 4180041812
    0xf9931f06: "13a_theuderic_00", // 4187168518
    0xf9c172da: "12_commanderpurcell", // 4190204634
    0xf9da6f4e: "0707_commander", // 4191842126
    0xfa0388dc: "99a_yonas_01", // 4194535644
    0xfa2eb797: "0410a_barkeeppolly", // 4197365655
    0xfa67049c: "0312_watcher", // 4201055388
    0xfa870c81: "12_eddayonkers", // 4203154561
    0xfae8037f: "9503_dog_s01_03", // 4209509247
    0xfb8b51cd: "13a_animalhunter_01", // 4220211661
    0xfbda4f40: "06_moirachapman", // 4225388352
    0xfbddebe1: "9505_skye_s01_07", // 4225625057
    0xfc29f87b: "9823_amelia", // 4230609019
    0xfcb02ff8: "9505_skye_s01_03", // 4239405048
    0xfcd5fecf: "1032_watcher", // 4241882831
    0xfcfb4504: "9502_zander_s01_04", // 4244325636
    0xfd152378: "9822_abayomi", // 4246020984
    0xfd85c766: "9503_dog_s01_07", // 4253402982
    0xfdd97187: "12_regentjinjae_01", // 4258886023
    0xfe324eb9: "0249a_aethelwynn_3", // 4264709817
    0xfe3f5424: "06_marceldupont03", // 4265563172
    0xfefd7b73: "1035_chef", // 4278025075
    0xff38e274: "0436a_tamraayad_01", // 4281918068
    0xff73e640: "01_gracemb_02", // 4285785664
    0xff79f28f: "12_ashedubois_01", // 4286182031
    0xff947fe4: "1659_anna", // 4287922148
    0xffee009d: "1624_imhotep_v4p3", // 4293787805
}
