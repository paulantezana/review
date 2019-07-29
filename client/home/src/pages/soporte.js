import React from "react"
import SiteLayout from "../layout/SiteLayout"
import QueueAnim from "rc-queue-anim"
import SEO from "../components/seo"
import { Button } from "antd";

export default ({ data }) => (
    <SiteLayout>
        <SEO
            title="El Mejor Soporte Técnico."
            description="Podrás pedir soporte y asesoría técnica por los diferentes medios que tenemos a continuación."
        />
        <div className="Center Container BannerB">
            <QueueAnim>
                <h1 key="h2">El Mejor Soporte Técnico.</h1>
                <p key="p">
                Podrás pedir soporte y asesoría técnica por los diferentes medios que tenemos a continuación.
                </p>
                <a href="https://dl.tvcdn.de/download/TeamViewer_Setup.exe" download>
                    <Button type="primary">Descarga (Soporte Remoto)</Button>
                </a>
            </QueueAnim>
        </div>
    </SiteLayout>
)
