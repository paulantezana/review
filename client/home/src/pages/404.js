import React from "react"
import { graphql } from "gatsby"

import SiteLayout from "../layout/SiteLayout"
import SEO from "../components/seo"
import { Result } from "antd";

class NotFoundPage extends React.Component {
    render() {
        const { data } = this.props
        const siteTitle = data.site.siteMetadata.title

        return (
            <SiteLayout location={this.props.location} title={siteTitle}>
                <SEO title="404: Not Found" />
                <Result
                    status="404"
                    title="404"
                    subTitle="Lo sentimos, la página que has visitado no existe."
                />
            </SiteLayout>
        )
    }
}

export default NotFoundPage

export const pageQuery = graphql`
    query {
        site {
            siteMetadata {
                title
            }
        }
    }
`
