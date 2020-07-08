import Head from 'next/head'
import Layout from '../components/layout'
import JoinableGameCard from '../components/joinable-game-card'

const Hall = () => {
  const joinableGames = [1, 2, 3, 4, 5, 6, 7, 8]
  return (
    <Layout>
      <Head>
        <title>52 Games are running!</title>
      </Head>
      <div className="container pt-32">
        <div className="flex justify-center">
          <div className="grid grid-cols-5 gap-8">
            {joinableGames.map((_, i) => ( 
              <div key={i}>
                <JoinableGameCard />
              </div>
            ))}
          </div>
        </div>
      </div>
    </Layout>
  )
}

export default Hall

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
