import React from "react"
import dashboard from "../../images/dashboard.jpg"
import slideHome from "../../images/slide-home.jpg"

// import { DataApp } from '../data/data';
import { Statistic } from "antd"

const { Countdown } = Statistic;
const deadline = new Date('2020','00','01').getTime();

export default () => (
    <div className="Slide">
        <img src={slideHome} className="Slide-bg" alt="slide-bg"/>
        <div className="Container Slide-data">
            <h1 className="Slide-title">Sistema Institucional Web</h1>
            <p className="Slide-text">Sistema web para gestionar los procesos de IEST</p>
            <Countdown 
                className="CountdownRelase"
                value={deadline}
                title={`Lanzamiento en: `}
                format="M [meses] - D [dias] - H [horas] - mm [minutos] - ss [segundos]"
            />
            <img
                src={dashboard}
                alt="Interfas del sistema"
                className="Slide-app"
            />
        </div>
    </div>
)
