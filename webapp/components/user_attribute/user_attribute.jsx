import React from 'react';
import {Tooltip, OverlayTrigger} from 'react-bootstrap';
import PropTypes from 'prop-types';

import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

export default class UserAttribute extends React.PureComponent {
    static propTypes = {
        theme: PropTypes.object.isRequired,
        gitlabURL: PropTypes.string.isRequired,
        username: PropTypes.string,
    };

    render() {
        const style = getStyle(this.props.theme);

        const username = this.props.username;

        if (!username) {
            return null;
        }

        return (
            <div style={style.container}>
            <a
                href={this.props.gitlabURL + '/' + username}
                target='_blank'
            >
                <i className='fa fa-gitlab'/>{' ' + username}
            </a>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            margin: '5px 0',
        },
    };
});
