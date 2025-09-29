import { View, StyleSheet } from 'react-native';
import AddMediaForm from './components/AddMedia/AddMediaForm';

export default function Index() {
  return (
    <View style={styles.container}>
      <AddMediaForm />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    padding: 16,
  },
});
