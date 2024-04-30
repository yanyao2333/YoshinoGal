import SideBar from "./sideBar";
import React from "react";

type Game = {
    name: string;
    poster_url: string;
    bg_url: string;
    description: string;
    details: Record<string, string>;
    screenshots: string[];
}

function GameDetail(props: { game: Game }) {
    return (
        <div className="flex h-screen overflow-hidden">
            <div className="flex flex-col gap-y-5 w-72 h-screen overflow-y-auto flex-shrink-0">
                <SideBar/>
            </div>
            <div className="flex-grow overflow-y-auto">
                {/* 背景 */}
                <div className="background-image h-64 bg-cover bg-repeat bg-center"
                     style={{backgroundImage: `url(${props.game.bg_url})`}}></div>
                {/* 游戏详情主体 */}
                <div
                    className="game-details mx-4 my-4 bg-white/60">
                    {/* 上半部分 */}
                    <div className="top-section flex justify-between items-center p-4 bg-gray-100 rounded-lg shadow">
                        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                            开始游戏
                        </button>
                        {/* 其他小按钮 */}
                        <button className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded">
                            在VNDB中查看
                        </button>
                    </div>
                    {/* 中部 */}
                    <div className=" divide-solid divide-y">
                        <div className="middle-section flex mt-4 p-4">
                            <img src={props.game.poster_url} alt="Game Poster" className="w-48 h-64 object-cover mr-8"/>
                            <div className="flex-1">
                                <h1 className="text-4xl font-bold mb-4">{props.game.name}</h1>
                                <div className="grid grid-cols-2 gap-y-8">
                                    {Object.entries(props.game.details).map(([key, value], index) => (
                                        <div key={index} className="flex flex-col gap-y-2">
                                            <span className="font-bold">{key}</span>
                                            <span>{value}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                        {/* 游戏简介 */}
                        <div className="description p-4">
                            <h1 className="text-3xl font-bold mb-2">简介</h1>
                            <p>{props.game.description}</p>
                        </div>
                        {/* 下半部分 */}
                        {/*<div className="bottom-section mt-4 p-4">*/}
                        {/*  <h1 className="text-3xl font-bold mb-2">游戏截图</h1>*/}
                        {/*  <div className="flex overflow-x-auto">*/}
                        {/*    {gameDetails.screenshots.map((screenshotUrl, index) => (*/}
                        {/*        <img key={index} src={screenshotUrl} alt={`Screenshot ${index + 1}`}*/}
                        {/*             className="screenshot mr-4 rounded-lg"*/}
                        {/*             style={{minWidth: "300px", height: "200px", objectFit: 'cover'}}*/}
                        {/*        />*/}
                        {/*    ))}*/}
                        {/*  </div>*/}
                        {/*</div>*/}
                        <div className="flex-grow overflow-y-auto">
                            <div className="bottom-section mt-4 p-4">
                                <h1 className="text-3xl font-bold mb-2">游戏截图</h1>
                                <div
                                    className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
                                    {props.game.screenshots.map((screenshotUrl, index) => (
                                        <img key={index} src={screenshotUrl} alt={`Screenshot ${index + 1}`}
                                             className="rounded-lg object-cover"
                                             style={{width: "100%", height: "200px"}}
                                        />
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default GameDetail;