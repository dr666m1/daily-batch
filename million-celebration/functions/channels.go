package celebrate

type channel struct {
	name      string
	playlists []string
	tags      []string
}

var channels = []channel{
	{
		name: "子兎音様",
		playlists: []string{
			"PL4PDyA42kQIz6SLG-2AuVP79xEpFGLoun", // 歌ってみた / Songs
			"PL4PDyA42kQIxNss1PEt99u76zQtemIgAO", // オリジナル曲
		},
		tags: []string{"天神子兎音"},
	},
	{
		name: "アイちゃん",
		playlists: []string{
			"PL0bHKk6wuUGIAmzzqdVMynRrAOi8odYFQ", // 歌ってみた・踊ってみた
			"PL0bHKk6wuUGLWGipKSf0dFrpuzDitERqD", // Kizuna AI Original Songs
		},
		tags: []string{"KizunaAI", "KizunaAIMusic"},
	},
	{
		name: "YuNiちゃん",
		playlists: []string{
			"PLjGAXih5gQmG6ZKg9QaNde_CvzKa2FwB5", // オリジナル
			"PLjGAXih5gQmHavTwLCBjCJRCi-a6Z4TG_", // 歌ってみた - YuNi
		},
		tags: []string{"YuNi"},
	},
	{
		name: "ちゃんまり",
		playlists: []string{
			"PLMSwJOUS0Tx4TGMWH8RN27xWyaD9gaoAd",        // cover
			"OLAK5uy_k5Yfyse6Dzm9lKu_SkwISRvwwJKjXF-2E", // from M
			"OLAK5uy_mMupmRhbWtZTGiP-RcDdVE5Ioilc4dy0M", // BiiRTHDAY
		},
		tags: []string{"かしこまり"},
	},
	{
		name: "葵ちゃん",
		playlists: []string{
			"PLFPYP7GcgzUQD-sMeyh_OO1DGIA4ez52A",        // 全カバー動画（投稿順）
			"OLAK5uy_mZmIsex6bXFJ2aaNtNX1CAi5dLffQbaII", // 有機的パレットシンドローム
			"OLAK5uy_kf0e0HNezVjnet8OUnHIK-5gEqkb8EVg4", // 声 ~Cover ch.~
		},
		tags: []string{"富士葵"},
	},
	{
		name: "ヒメヒナ",
		playlists: []string{
			"PL1tX8zAv8bPkOg13XRuyrNsFQw8X2d6CQ", // ヒメヒナMusic
		},
		tags: []string{"ヒメヒナ"},
	},
}
