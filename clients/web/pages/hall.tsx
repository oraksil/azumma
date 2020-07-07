import React from 'react'
import Layout from '../components/layout'
import Head from 'next/head'

type Props = {
}

const Hall = () => {
  return (
    <React.Fragment>
      <Layout>
        <Head>
          <title>52 Games are running!</title>
        </Head>
        <h1>Hall Page</h1>
      </Layout>
    </React.Fragment>
  )
}

export default Hall

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
