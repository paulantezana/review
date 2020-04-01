import React, { Component, Fragment } from 'react';
import { Button, Divider, Tooltip,  Modal, Table } from 'antd';

import FormModal from './Form';
import ToolBar from '../../components/ToolBar';
import { institutionPaginate, institutionCreate, institutionUpdate } from '../../../services/institution';

class DataList extends Component {
    state = {
        data: {
            list: [],
            pagination: {},
        },
        loading: false,

        modalVisible: false,
        currentItem: {},
        modalType: 'create',
    }

    componentDidMount() {
        this.onQueryAll();
    }

    onQueryAll = (param = {}) => {

        const params = {
            ...param,
            limit: param.limit,
        }

        institutionPaginate(params).then(response => {
            if(response.success){
                this.setState({
                    data: {
                        list: response.data,
                        pagination: {
                            current: response.current_page,
                            total: response.total,
                            pageSize: response.limit,
                        },
                    }
                })
            } else {
                Modal.error({ title: 'Institución', content: response.message });
            }
        });
    }

    onShowModal = (modalType, currentItem = {}) => {
        this.setState({
            modalVisible: true,
            currentItem: currentItem,
            modalType: modalType,
        })
    };

    onSaveModal = values => {
        const modalType = this.state.modalType;
        if(modalType === 'create'){
            institutionCreate(values).then(response => {
                if(response.success){
                    this.onQueryAll();
                }else{
                    Modal.error({ title: 'Institución', content: response.message });
                }
            });
        }else if(modalType === 'update'){
            institutionUpdate(values).then(response => {
                if(response.success){
                    this.onQueryAll();
                }else{
                    Modal.error({ title: 'Institución', content: response.message });
                }
            });
        }
    }

    onCancelModal = () => {
        this.setState({
            modalVisible: false,
            currentItem: {},
            modalType: 'create',
            loading: false,
        })
    }

    render() {
        const { loading } = this.state;

        const columns = [
            {
                title: 'Institución',
                dataIndex: 'institute',
                key: 'institute',
            },
            {
                title: 'Denominación',
                dataIndex: 'prefix',
                key: 'prefix',
            },
            {
                title: 'Denominación',
                dataIndex: 'prefix_short_name',
                key: 'prefix_short_name',
            },
            {
                title: 'Director',
                dataIndex: 'director',
                key: 'director',
            },
        ];

        const paginationProps = {
            showSizeChanger: true,
            showQuickJumper: true,
            ...this.state.data.pagination,
          };

        return (
            <Fragment>
                <ToolBar
                    title="Programas de estudio"
                    rightActions={[
                        <Tooltip title="Nuevo">
                            <Button
                                type="primary"
                                loading={loading}
                                onClick={() => this.onShowModal('create')}
                                icon="plus"
                            />
                        </Tooltip>,
                        <Divider type="vertical" />,
                        <Tooltip title="Regrescar">
                            <Button
                                icon="reload"
                                loading={loading}
                                onClick={this.onQueryAll}
                            />
                        </Tooltip>
                    ]}
                />
                <Table
                    columns={columns}
                    size="small"
                    dataSource={this.state.data.list} 
                    // pagination={paginationProps} 
                />
                <FormModal 
                    loading={this.state.loading}
                    modalVisible={this.state.modalVisible}
                    currentItem={this.state.currentItem}
                    onSave={this.onSaveModal}
                    onCancel={this.onCancelModal}
                />
            </Fragment>
        );
    }
}

export default DataList;
