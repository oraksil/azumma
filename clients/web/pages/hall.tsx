import React from 'react'
import Container from '../components/container'
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
        <Container>
          <h1>Hall Page</h1>
        </Container>
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
