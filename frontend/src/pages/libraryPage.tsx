import SideBar from "../components/sideBar";
import GamePosterWall from "../components/gamePosterWall";

async function getPosterWall() {
    const backend_port = 8080;
    const game_library_path = "E:\\GalGames";
    const res = await fetch(`http://localhost:` + backend_port + `/library/index/posterwall`, {
        method: 'POST',
        body: JSON.stringify({
            directory: game_library_path
        }),
        cache: 'no-store',
    });
    const data = await res.json();

    if (!data) {
        return {
            notFound: true,
        }
    }

    if (data["code"] !== 0) {
        return {
            notFound: true,
        }
    }

    // games是个字典，key是游戏名，value是海报路径
    let games = data["data"];

    let gamesShowList = [];

    let idx = 1;

    for (let game in games) {
        gamesShowList.push({
            name: game,
            imageSrc: 'http://localhost:' + backend_port + '/img/' + games[game],
            id: idx,
            href: '/game/' + game,
            imageAlt: game,
        });
    }

    console.log(gamesShowList)

    return gamesShowList
}


export default async function Home() {
    const gamesShowList = await getPosterWall();
    // const gamesShowList = example;
    return (
        <div className="flex bg-gray-100">
            <div className="flex flex-col gap-y-5 w-72 h-screen overflow-y-auto">
                <SideBar/>
            </div>
            <div className="flex-grow overflow-y-auto h-screen">
                <GamePosterWall game={gamesShowList}/>
            </div>
        </div>
    );

}