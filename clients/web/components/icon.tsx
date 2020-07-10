import dynamic from 'next/dynamic'

type Props = {
  name: string,
  className: string,
}

const Icon = (props: Props) => {
  // this is a workaround to import module by variable
  // https://github.com/webpack/webpack/issues/6680#issuecomment-370800037
  const IconFromSvg = dynamic(() => import('../public/assets/icons/' + props.name + '.svg'))
  return (
    <div className={props.className}>
      <IconFromSvg />
    </div>
  )
}

export default Icon 
