import React, { Fragment } from 'react';
import { Icon, Divider } from 'antd';

const ToolBar = ({className, leftActions = [], rightActions = [], title = "", subTitle = "" })=>{
    return (
        <div className={`SnToolbar ${className}`} >
            <div className='SnToolbar-left'>
                <Icon type="bars" />
                {
                    (title != "")  && (
                        <Fragment>
                            <Divider type="vertical"/>
                            <strong>{title}</strong> 
                        </Fragment> 
                    )
                }
                {
                    (subTitle != "")  && (
                        <Fragment>
                            <Divider type="vertical"/>
                            {subTitle}
                        </Fragment> 
                    )
                }
                { 
                    leftActions.map(item=>(
                        <Fragment>
                            { item }
                        </Fragment>
                    )) 
                }
            </div>
            <div className='SnToolbar-right'>
                { 
                    rightActions.map((item, index)=>(
                        <Fragment key={index}>
                            { item }
                        </Fragment>
                    )) 
                }
            </div>
        </div>
    )
}

export default ToolBar;