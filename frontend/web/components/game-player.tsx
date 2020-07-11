import styles from './game-player.module.css'
import Icon from './icon'

const GamePlayer = () => {
  return (
    <div className={styles['player-frame']}>
      <div className={styles['player']}>
        <video className="w-full h-full bg-purple-400" autoPlay={true}> 
          <source src="http://techslides.com/demos/sample-videos/small.webm" type="video/webm" /> 
        </video>
      </div>
      <div className={styles['player-ctl']}>
        <div className="absolute left-0 top-0 h-full px-2">
          <div className="inline-block p-2">
            <Icon name="toll" className="fill-current text-gray-200" />
          </div>
          <div className="inline-block p-2">
            <Icon name="gamepad" className="fill-current text-gray-200" />
          </div>
        </div>
        <div className="flex justify-center h-full"> 
          <div className="inline-block p-2">
            <Icon name="link" className="fill-current text-gray-200" />
          </div>
        </div>
        <div className="absolute right-0 top-0 h-full px-2">
          <div className="inline-block p-2">
            <Icon name="pause" className="fill-current text-gray-200" />
          </div>
          <div className="inline-block p-2">
            <Icon name="volume-off" className="fill-current text-gray-200" />
          </div>
        </div>
      </div>
    </div>
  )
}

export default GamePlayer
