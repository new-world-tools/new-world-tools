package azcs

var hashStatusEffectDataEffectCategories = map[uint32]string{
    0x00000000: "", // 0
    0x01970eea: "frost", // 26676970
    0x059278a3: "dot", // 93485219
    0x06a98af5: "tokengatherboost", // 111774453
    0x0925b47e: "mutatorcurse", // 153465982
    0x0bf0d05e: "unclampedexhaust", // 200331358
    0x0c29ef40: "armorfortify", // 204074816
    0x0e6c9dbe: "silence", // 241999294
    0x0e6eec05: "essencerupturet3", // 242150405
    0x0ea6e98c: "fireburn", // 245819788
    0x0f3b6ac1: "disease", // 255552193
    0x11756a55: "rainofarrows", // 292907605
    0x12271997: "harvestmana", // 304552343
    0x153e9a58: "anointed", // 356424280
    0x16a1723b: "nondispellabledebuff", // 379679291
    0x16f4f95b: "root", // 385153371
    0x1a237df1: "housingdmgbeasts", // 438533617
    0x1fcf5da8: "rend", // 533683624
    0x21f8847e: "resistboost", // 569934974
    0x253efcc2: "throwingaxe", // 624884930
    0x271689ad: "debuff", // 655788461
    0x2c975ce7: "threat", // 748117223
    0x319b9e70: "attributes", // 832282224
    0x350771dd: "slow", // 889680349
    0x356410d3: "unclampedfortify", // 895750355
    0x35cbaa20: "clearonreset", // 902539808
    0x37d2b8b0: "nullchamberbuff", // 936556720
    0x3b059b76: "traproot", // 990223222
    0x3b8ccc74: "stamina", // 999083124
    0x3dfe1fb9: "harvesthp", // 1040064441
    0x3e3f641a: "sturgeonfishingbuff", // 1044341786
    0x3ed8a605: "insmoke", // 1054385669
    0x3f83eb72: "voidgauntlet", // 1065610098
    0x412aed17: "overheatdebuff", // 1093332247
    0x41f7d51e: "gsoffenseonly", // 1106760990
    0x45c36aca: "immortal", // 1170434762
    0x46700042: "infectedthrow", // 1181745218
    0x4691f41a: "freeze", // 1183970330
    0x4869a8e3: "bleed", // 1214884067
    0x4a198fba: "sprintdelaycombat", // 1243189178
    0x4df4d6a1: "resisttincture", // 1307891361
    0x5074594e: "greatswordbuff", // 1349802318
    0x50e366cd: "uncleansabledot", // 1357080269
    0x51f9be47: "shock", // 1375321671
    0x55521d04: "burn", // 1431444740
    0x56642768: "boost", // 1449404264
    0x56ff898d: "mutatorhearty", // 1459587469
    0x5838412a: "exhaustedstate", // 1480081706
    0x5f0b8b89: "slow&stun", // 1594592137
    0x61e3e397: "tradeskill", // 1642324887
    0x62c04ae9: "focus", // 1656769257
    0x64d9d4e5: "huntervision", // 1691997413
    0x66a32e6b: "carryboulder", // 1721970283
    0x6795d70d: "fishingbuff", // 1737873165
    0x6910ce50: "mutatorempower", // 1762709072
    0x6a23fc9c: "stoneform", // 1780743324
    0x6ad0de1a: "timer", // 1792073242
    0x6b4075c6: "iceroot", // 1799386566
    0x7348d0bc: "unclampedfactionresourcemod", // 1934151868
    0x75274263: "buff", // 1965507171
    0x75516071: "passivefoodregen", // 1968267377
    0x7969dc93: "essencerupturet2", // 2036980883
    0x7a4af000: "warhammer", // 2051731456
    0x7bbd9903: "frostbite", // 2076023043
    0x7fd16d73: "admiralbruteresistdispell", // 2144431475
    0x8045b257: "teamwipecat", // 2152051287
    0x81518a58: "foodinitialrecovery", // 2169604696
    0x82fd433c: "stance", // 2197635900
    0x83af03fd: "weaken", // 2209285117
    0x847dd78f: "nondispellableempower", // 2222839695
    0x880baad1: "exhaust", // 2282466001
    0x8cc0104e: "unclampeddisease", // 2361397326
    0x8d2e6b33: "gatherboost", // 2368629555
    0x8dbf1fb2: "unclampedrend", // 2378112946
    0x8f7ce27e: "frostbuff", // 2407326334
    0x90ea6d55: "ignoredebuffs", // 2431282517
    0x9145c383: "debilitate", // 2437268355
    0x934dd900: "ccdurreduction", // 2471352576
    0x963ebda8: "luck", // 2520694184
    0x9821eb23: "ai_neutral", // 2552359715
    0x98c692d4: "bblastchance", // 2563150548
    0x9a65fe8e: "flamekeeperdebuff", // 2590375566
    0x9d205605: "empower", // 2636142085
    0x9f9e8ea3: "spear", // 2677968547
    0xa48c6377: "deathfog", // 2760663927
    0xa5e6215b: "family", // 2783322459
    0xa6e62e58: "defboost", // 2800103000
    0xa77733c7: "unclampedslow", // 2809607111
    0xa9287eae: "heal", // 2838003374
    0xac632aaf: "tracker", // 2892180143
    0xaf283bd4: "unclampedazothmod", // 2938649556
    0xafb7acf1: "laststand", // 2948050161
    0xb218f86b: "nondispellableweaken", // 2987980907
    0xb219dfe7: "affliction", // 2988040167
    0xb2d28d9b: "hearty", // 3000143259
    0xb31a5bff: "dmg", // 3004849151
    0xb69f6a5c: "fortify", // 3063900764
    0xb6d04d9b: "regen", // 3067104667
    0xbcb318ff: "siegeonly", // 3165853951
    0xbcf55e04: "crit", // 3170196996
    0xc92581e4: "azothstaff", // 3374678500
    0xcacc3a47: "foodutility", // 3402381895
    0xccc6a4db: "haste", // 3435570395
    0xccf769b5: "energized", // 3438766517
    0xcee75d23: "traversalhaste", // 3471269155
    0xd11b68a6: "essencerupture", // 3508234406
    0xd27bd9ee: "void", // 3531332078
    0xd43829f7: "food", // 3560450551
    0xd5d022e2: "showerroot", // 3587187426
    0xd604308a: "nondispellablebuff", // 3590598794
    0xd67a656d: "admiralbrutetransitiondispell", // 3598345581
    0xd71fe96f: "musicbuff", // 3609192815
    0xdbb21a79: "cc", // 3685882489
    0xdcc98685: "bossdebuff", // 3704194693
    0xe048bcba: "powderburn", // 3762863290
    0xe0608d29: "essencerupturet1", // 3764423977
    0xe493e3cb: "uninterruptible", // 3834897355
    0xe702b4b1: "carrycannonball", // 3875714225
    0xed3eb615: "armorrend", // 3980310037
    0xf120340f: "lifestaffbuff", // 4045419535
    0xf1cf77b4: "cleanse", // 4056905652
    0xf34393c9: "recovery", // 4081292233
    0xf374d10e: "stun", // 4084519182
    0xf9741949: "poison", // 4185135433
    0xfb0bf294: "mutatordisease", // 4211864212
    0xfebc9e4f: "trapper", // 4273774159
}