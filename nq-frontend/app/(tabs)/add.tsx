import React from "react";
import AddMediaPage from "../../src/pages/AddMediaPage";
import { useApolloClient } from "@apollo/client/react";
import { createMedia } from "../../lib/createMedia";
import { createActivity } from "../../lib/createActivity";
import { GET_HOME_MEDIA_QUERY, ME_ACTIVITIES_QUERY } from "../../lib/graphql";
import { Media } from "../../src/types";
import { useAuth } from "../../lib/AuthContext";

export default function AddTabPage() {
  const apolloClient = useApolloClient();
  const { hasToken } = useAuth();
  const [isAddingMedia, setIsAddingMedia] = React.useState(false);

  const handleAddMedia = async (newMedia: Omit<Media, "id">) => {
    setIsAddingMedia(true);
    try {
      const result = await createMedia(newMedia);
      if (result?.id) {
        await createActivity({
          mediaId: result.id,
          statusId: 1,
        });
      }
      const queries: Promise<unknown>[] = [
        apolloClient.query({ query: GET_HOME_MEDIA_QUERY, fetchPolicy: "network-only" }),
      ];
      if (hasToken) {
        queries.push(apolloClient.query({ query: ME_ACTIVITIES_QUERY, fetchPolicy: "network-only" }));
      }
      await Promise.all(queries);
    } catch (error) {
      console.error("Failed to add media:", error);
    } finally {
      setIsAddingMedia(false);
    }
  };

  return <AddMediaPage onBack={() => {}} onAddMedia={handleAddMedia} isLoading={isAddingMedia} />;
}
