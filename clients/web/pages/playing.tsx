import React from 'react'
import Container from '../components/container'
import Layout from '../components/layout'
import Head from 'next/head'

type Props = {
}

const Playing = () => {
  return (
    <React.Fragment>
      <Layout>
        <Head>
          <title>Street Fighter II</title>
        </Head>
        <Container>
          <h1>Playing Page</h1>
        </Container>
      </Layout>
    </React.Fragment>
  )
}

export default Playing

export const getStaticProps = async () => {
  return {
    props: {},
  }
}
