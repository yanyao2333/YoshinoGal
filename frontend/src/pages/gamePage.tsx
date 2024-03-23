async function getGameDetail(game_name: string) {
    const backend_port = 8080;
    const game_library_path = "E:\\GalGames";
    const metadata_res = await fetch(`http://localhost:` + backend_port + `/library/metadata/get/one`, {
        method: 'POST',
        body: JSON.stringify({
            game_name: game_name,
            directory: game_library_path
        }),
        cache: 'no-store',
    });
    const meta = await metadata_res.json();
    if (!meta) {
        return {
            notFound: true,
        }
    }

    if (meta["code"] !== 0) {
        return {
            notFound: true,
        }
    }
    const game_path = meta["data"]["local_poster_path"].replace(/\/metadata\/.*/, "");
    const playtime_res = await fetch(`http://localhost:` + backend_port + `/playtime/get/one`, {
        method: 'POST',
        body: JSON.stringify({
            directory: game_path,
        }),
        cache: 'no-store',
    });
    const playtime_data = await playtime_res.json();
    console.log(playtime_data)

    if (!playtime_data) {
        return {
            notFound: true,
        }
    }

    if (playtime_data["code"] !== 0) {
        return {
            notFound: true,
        }
    }

    let gameDetail = meta["data"];

    let random_bg = gameDetail["local_screenshots_paths"][Math.floor((Math.random() * gameDetail["local_screenshots_paths"].length))]

    if (gameDetail["length_minutes"] === 0) {
        switch (gameDetail["length"]) {
            case 1:
                gameDetail["length_minutes"] = "非常短";
                break;
            case 2:
                gameDetail["length_minutes"] = "短";
                break;
            case 3:
                gameDetail["length_minutes"] = "中等";
                break;
            case 4:
                gameDetail["length_minutes"] = "长";
                break;
            case 5:
                gameDetail["length_minutes"] = "非常长";
                break;
        }
    }

    switch (gameDetail["dev_status"]) {
        case 0:
            gameDetail["dev_status"] = "已完成";
            break;
        case 1:
            gameDetail["dev_status"] = "正在开发";
            break;
        case 2:
            gameDetail["dev_status"] = "已取消";
            break;
    }

    let processedGameDetail = {
        name: gameDetail["name"],
        bg_url: 'http://localhost:' + backend_port + '/img/' + random_bg,
        description: gameDetail["description"],
        details: {
            "发行时间": gameDetail["release_date"],
            "VNDB评分": gameDetail["score"]["VNDB"],
            "预计游玩时间": String(gameDetail["length_minutes"]) + "min",
            "开发状态": gameDetail["dev_status"],
            "当前已玩": String(Math.round(playtime_data["data"])) + "min",
            // "当前已玩": "316min"
        },
        poster_url: 'http://localhost:' + backend_port + '/img/' + gameDetail["local_poster_path"],
        screenshots: gameDetail["local_screenshots_paths"].map((screenshot: string) => {
            return 'http://localhost:' + backend_port + '/img/' + screenshot;
        })
    };

    return processedGameDetail
}

export default async function GamePage({params}: { params: { game_name: string } }) {
    const game_name = params.game_name
    if (game_name === null) {
        return {
            notFound: true,
        }
    }
    if (!game_name) {
        return {
            notFound: true,
        }
    }
    const gameDetail = await getGameDetail(decodeURI(game_name))
    console.log(gameDetail)

    return (
        <div>
            <GamePage game={gameDetail}/>
        </div>
    )
}