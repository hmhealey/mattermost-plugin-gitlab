import PropTypes from 'prop-types';
import React from 'react';
import {Tooltip, OverlayTrigger} from 'react-bootstrap';

export default class SidebarButton extends React.PureComponent {
    static propTypes = {
        count: PropTypes.number,
        href: PropTypes.string,
        icon: PropTypes.node.isRequired,
        onClick: PropTypes.func,
        style: PropTypes.object,
        tooltipId: PropTypes.string,
        tooltipPlacement: PropTypes.string,
        tooltipText: PropTypes.string,
    };

    render() {
        let href = this.props.href;
        let target = '_blank';
        let rel = 'noopener noreferrer';
        if (href == null && this.props.onClick) {
            href = '#';
            target = '';
            rel = '';
        }

        let buttonText;
        if (this.props.count != null) {
            buttonText = ' ' + this.props.count;
        }

        let button = (
            <a
                href={href}
                onClick={this.props.onClick}
                target={target}
                rel={rel}
                style={this.props.style}
            >
                {this.props.icon}
                {buttonText}
            </a>
        );

        if (this.props.tooltipText) {
            button = (
                <OverlayTrigger
                    placement={this.props.tooltipPlacement}
                    overlay={<Tooltip id={this.props.tooltipId}>{this.props.tooltipText}</Tooltip>}
                >
                    {button}
                </OverlayTrigger>
            );
        }

        return button;
    }
}
