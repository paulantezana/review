import React from "react"
import { graphql } from "gatsby"
import ApiLayout from "../layout/ApiLayout"
import SEO from "../components/seo"
import { Anchor } from "antd"

const GetAnchors = ({ headings }) => {
    const anchors = []
    headings.forEach(item => {
        if (item.depth === 2) {
            const al = item.value.replace(" ", "-")
            const linkA = `#${al.toLowerCase()}`
            anchors.push(
                <Anchor.Link
                    key={linkA.trim()}
                    href={linkA.trim()}
                    title={item.value}
                />
            )
        }
    })
    return anchors
}

export default props => {
    const post = props.data.markdownRemark
    const siteTitle = props.data.site.siteMetadata.title
    return (
        <ApiLayout location={props.location} title={siteTitle}>
            <SEO title={post.frontmatter.title} />
            <article className="ApiPage">

                <Anchor className="ApiPage-anchor">
                    <GetAnchors headings={post.headings} />
                </Anchor>
                <div className="ApiPage-content">
                    <h1>{post.frontmatter.title}</h1>
                    <div
                        className="Markdown"
                        dangerouslySetInnerHTML={{ __html: post.html }}
                    />
                </div>
            </article>
        </ApiLayout>
    )
}

export const query = graphql`
    query($slug: String!) {
        site {
            siteMetadata {
                title
                author
            }
        }
        markdownRemark(fields: { slug: { eq: $slug } }) {
            html
            frontmatter {
                title
            }
            headings {
                value
                depth
            }
        }
    }
`
