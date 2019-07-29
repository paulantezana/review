import React from "react"
import SiteLayout from "../layout/SiteLayout"
import QueueAnim from "rc-queue-anim"
import SEO from "../components/seo"

export default ({ data }) => (
    <SiteLayout>
        <SEO
            title="Términos de licencia."
            description="Términos de licencias de sistemas para IESTP."
        />
        <div className="Center Container BannerB">
            <QueueAnim>
                <h1 key="h2">Términos de licencia.</h1>
                <p key="p">
                El sistema fue desarrollado por los estudiantes del ISTP VILCANOTA promoción 2018.
                    La mejor promoción durante toda la historia del programa de estudios computación e informática.
                </p>
            </QueueAnim>
        </div>
    </SiteLayout>
)
