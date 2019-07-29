import React from "react"
import dashboard from "../images/dashboard.jpg"
import slideHome from "../images/slide-home.jpg"

import { DataApp } from '../data/data';

export default () => (
    <div className="Slide">
        <img src={slideHome} className="Slide-bg" alt="slide-bg"/>
        <div className="Container Slide-data">
            <h1 className="Slide-title">Sistema institucional</h1>
            <p className="Slide-text">Sistema web para gestionar los procesos de IEST.</p>
            <p className="Slide-version">{DataApp.version}</p>
            <img
                src={dashboard}
                alt="Interfas del sistema"
                className="MainBanner_img"
            />
        </div>
    </div>
)
