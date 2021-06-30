import styled from 'styled-components';
import { space, color, typography, border, shadow } from 'styled-system';

const Textarea = styled.textarea`
border: 1px solid ${props => props.theme.colors.black};
outline: none;
appearance: none;
width: 100%;
resize: vertical;
caret-color: ${props => props.theme.colors.primary};

transition: ${props => props.theme.transitions[0]};

&:focus {
  border-color: ${props => props.theme.colors.primary};
}
${space} ${color} ${typography} ${border} ${shadow}
`;

Textarea.defaultProps = {
  color: 'white',
  bg: 'grays.0',
  borderRadius: 1,
  padding: 3,
  fontSize: 2,
  fontWeight: 0,
  border: 0,
  borderColor: 'white',
};

export default Textarea;
