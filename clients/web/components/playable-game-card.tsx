const PlayableGameCard = () => {
  return (
    <div className="card-frame inline-block w-48 border-solid border-gray-200 border-2">
      <img className="cover-image h-48 object-center object-cover"
        src="https://tailwindcss.com/img/tailwind-ui-sidebar.png" />
      <div className="p-2">
        <div className="game-title text-sm">
          <span className="font-bold">Game Title</span>
        </div>
        <div className="game-producer text-xs">Game Producer</div>
        <div className="game-players">
          <span className="text-sm">2</span>
          <span className="text-xs">P</span>
        </div>
      </div>
    </div>
  )
}

export default PlayableGameCard 
