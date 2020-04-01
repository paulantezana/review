import React from 'react';

import { getAppModule } from '../../services/public';
import { Modal, List } from 'antd';

class Module extends React.Component{
    state = {
        data: {
            list: [],
            pagination: {},
        },
        loading: false,
    }

    componentDidMount(){
        this.setState({
            loading: true,
        });
        getAppModule({
            limit: 100
        }).then(response => {
            if(response.success){
                this.setState({
                    data: {
                        list: response.data,
                        pagination: {
                            current: response.current_page,
                            total: response.total,
                            pageSize: response.limit,
                        },
                    },
                    
                })
            } else {
                Modal.error({ title: 'Institución', content: response.message });
            }
        }).finally(e=>{
            this.setState({
                loading: false,
                modalVisible: false,
            });
        });
    }
    render(){
        return (
            <List
                grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 4 }}
                dataSource={this.state.data.list}
                renderItem={item => (
                    <List.Item>
                        <div>
                            <h3>{ item.name }</h3>
                            <p>{ item.description }</p>
                            {/* <Link to={item.doc}>Seguir leendo</Link> */}
                        </div>
                    </List.Item>
                )}
            />
        )
    }
}

export default Module;